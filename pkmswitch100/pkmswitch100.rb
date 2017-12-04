#!/usr/bin/env ruby

def optParser
  opts = {}
  ARGV.each do |arg|
    arg_arr = arg.split('=')
    opts[arg_arr[0].delete('-')] = arg_arr[1]
  end
  opts
end

def vercomp(old = '', new = '')
  # 3.2.2-3.4 3.2.2-3.3 3.2.2-1.1
  ov = old.gsub(/-.*$/, '')
  nv = new.gsub(/-.*$/, '')
  orel = old.gsub(/^.*-/, '').to_f
  nrel = new.gsub(/^.*-/, '').to_f
  # compare major version
  if nv.to_f > ov.to_f
    return 1
  elsif nv.to_f < ov.to_f
    return -1
  else
    # if the first 2 digits are the same, compare the last digit to make sure
    ovn = ov.gsub(/^.*\./, '').to_i
    nvn = nv.gsub(/^.*\./, '').to_i
    if nvn > ovn
      return 1
    elsif nvn < ovn
      return -1
    else
      # need to compare release numbers
      if nrel > orel
        return 1
      elsif nrel < orel
        return -1
      else
        return 0
      end
    end
  end
end

def max(arr = [])
  size = arr.length
  max = arr[size - 1]
  size.times do |i|
    max = arr[i] if vercomp(arr[i], max) < 0
  end
  max
end

def min(arr = [])
  size = arr.length
  min = arr[size - 1]
  size.times do |i|
    min = arr[i] if arr[i] < min
  end
  min
end

ffmpeg3_arr = %w(ffmpeg libavcodec57 libavdevice57 libavfilter6 libavformat57 libavresample3 libavutil55 libpostproc54 libswresample2 libswscale4)
vlc_arr = ['libvlc5', 'libvlccore8', 'vlc', 'vlc-codec-gstreamer', 'vlc-noX', 'vlc-qt', 'vlc-codecs', 'npapi-vlc']
gstreamer_arr = ['gstreamer-plugins-bad', 'gstreamer-plugins-libav', 'gstreamer-plugins-ugly', 'gstreamer-plugins-ugly-orig-addon', 'libgstadaptivedemux-1_0-0', 'libgstbadaudio-1_0-0', 'libgstbadbase-1_0-0', 'libgstbadvideo-1_0-0', 'libgstbasecamerabinsrc-1_0-0', 'libgstcodecparsers-1_0-0', 'libgstgl-1_0-0', 'libgstmpegts-1_0-0', 'libgstphotography-1_0-0', 'libgsturidownloader-1_0-0', 'libgstwayland-1_0-0', 'gstreamer-plugins-bad-orig-addon', 'libgstinsertbin-1_0-0', 'libgstplayer-1_0-0', 'libgstvdpau']
all_arr = ffmpeg3_arr + vlc_arr + gstreamer_arr

opts = optParser
pkgs = {}
all_arr = ffmpeg3_arr if opts.key?('ffmpeg')
all_arr = vlc_arr if opts.key?('vlc')
all_arr = gstreamer_arr if opts.key?('gstreamer')

pool_size = all_arr.length + 1
jobs = Queue.new
all_arr.each { |i| jobs.push i }

workers = Array.new(pool_size) do
  Thread.new do
    begin
      while x = jobs.pop(true)
        stdin = `LANG=en_US.UTF-8 zypper --no-refresh se -v #{x}`
        stdin.each_line do |line|
          next unless line =~ /^(i(\+)?|v)\s+\|/ && line.index(x + "\s")
          # v | libavformat57   | package | 3.2.2-3.4 | i586   | packman
          arr = line.split('|').each(&:strip!)
          arr[0] = "i" if arr[0] == 'i+'
          if pkgs.key?(arr[1])
            if pkgs[arr[1]].key?(arr[0])
              if pkgs[arr[1]][arr[0]].key?(arr[3])
                pkgs[arr[1]][arr[0]][arr[3]] << arr[4]
              else
                pkgs[arr[1]][arr[0]][arr[3]] = [arr[4]]
              end
            else
              pkgs[arr[1]][arr[0]] = { arr[3] => [arr[4]] }
            end
          else
            pkgs[arr[1]] = { arr[0] => { arr[3] => [arr[4]] } }
          end
        end
      end
    rescue ThreadError
    end
  end
end

workers.map(&:join)

oss = []
update = []
pkgs.each do |k, v|
  instdver = v['i'].keys[0]
  pkmver = max(v['v'].keys).to_s
  if v['v'].keys.include?(instdver)
    # not from packman
    oss << k if instdver != pkmver
  else
    # if the installed package's release number is near to packman's
    # then you have outdated packman version. if near to oss's
    # then it is oss package.
    instdrel = instdver.gsub(/^.*-/, '')
    vers = v['v'].keys.map { |i| i.gsub(/^.*-/, '') }
    diff = vers.map { |j| j.to_f - instdrel.to_f }.map do |m|
      m < 0 ? m * -1 : m
    end
    mini = min(diff)
    pos = diff.find_index(mini)
    near = v['v'].keys[pos]
    if near == pkmver
      update << k
    else
      oss << k
    end
  end
end

puts '======================= Packages not from Packman ========================='

if oss.empty?
  puts 'Good! All packages are from Packman!'
else
  str = ''
  oss.each do |i|
    puts i
    # zypper install libavformat57-3.2.2-3.4.x86_64
    pkm = max(pkgs[i]['v'].keys).to_s
    arch = pkgs[i]['i'].values[0][0]
    full = i + '-' + pkm + '.' + arch
    str = str + ' ' + full
  end
  puts "\nFIX: Run 'sudo zypper install" + str + "'."
end

puts '======================= Packman Packages need updates ====================='

if update.empty?
  puts 'Good! All packages are from Packman and at their latest versions!'
else
  str = ''
  update.each do |i|
    puts i
    str = str + ' ' + i
  end
  puts "\nFIX: Run 'sudo zypper up" + str + "'."
end
