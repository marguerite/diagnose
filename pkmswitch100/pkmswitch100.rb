#!/usr/bin/env ruby

def optParser()
	opts = {}
	ARGV.each do |arg|
		arg_arr = arg.split("=")
		opts[arg_arr[0].gsub("-","")] = arg_arr[1]
	end
	return opts
end

def vercomp(old="",new="")
	# 3.2.2-3.4 3.2.2-3.3 3.2.2-1.1 
	ov,nv = old.gsub(/-.*$/,''),new.gsub(/-.*$/,'')
	orel,nrel = old.gsub(/^.*-/,'').to_f,new.gsub(/^.*-/,'').to_f
	# compare major version
	if nv.to_f > ov.to_f
		return 1
	elsif nv.to_f < ov.to_f
		return -1
	else
	        # if the first 2 digits are the same, compare the last digit to make sure
		ovn,nvn = ov.gsub(/^.*\./,'').to_i,nv.gsub(/^.*\./,'').to_i
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

def max(arr=[])
	size = arr.length	
	max = arr[size-1]
	size.times do |i|
		if vercomp(arr[i],max) < 0
			max = arr[i]
		end
	end
	return max
end

def min(arr=[])
	size = arr.length
	min = arr[size-1]
	size.times do |i|
		if arr[i] < min
			min = arr[i]
		end
	end
	return min
end

ffmpeg3_arr = ["ffmpeg","libavcodec57","libavdevice57","libavfilter6","libavformat57","libavresample3","libavutil55","libpostproc54","libswresample2","libswscale4"]
vlc_arr = ["libvlc5","libvlccore8","vlc","vlc-codec-gstreamer","vlc-noX","vlc-qt","vlc-codecs","npapi-vlc"]
gstreamer_arr = ['gstreamer-plugins-bad','gstreamer-plugins-libav','gstreamer-plugins-ugly','gstreamer-plugins-ugly-orig-addon','libgstadaptivedemux-1_0-0','libgstbadaudio-1_0-0','libgstbadbase-1_0-0','libgstbadvideo-1_0-0','libgstbasecamerabinsrc-1_0-0','libgstcodecparsers-1_0-0','libgstgl-1_0-0','libgstmpegts-1_0-0','libgstphotography-1_0-0','libgsturidownloader-1_0-0','libgstwayland-1_0-0','gstreamer-plugins-bad-orig-addon','libgstinsertbin-1_0-0','libgstplayer-1_0-0','libgstvdpau']
all_arr = ffmpeg3_arr + vlc_arr + gstreamer_arr

opts,pkgs = optParser,{}
all_arr = ffmpeg3_arr if opts.has_key?("ffmpeg")
all_arr = vlc_arr if opts.has_key?("vlc")
all_arr = gstreamer_arr if opts.has_key?("gstreamer")

pool_size = all_arr.length + 1
jobs = Queue.new
all_arr.each {|i| jobs.push i}

workers = pool_size.times.map do
	Thread.new do
		begin
			while x = jobs.pop(true)
				stdin = `LANG=en_US.UTF-8 zypper --no-refresh se -v #{x}`
				stdin.each_line do |line|
					if (line.index("v\s|") || line.index("i\s|")) && line.index(x + "\s")
						# v | libavformat57   | package | 3.2.2-3.4 | i586   | packman
						arr = line.split("|").each {|i| i.strip!}
						unless pkgs.has_key?(arr[1])
							pkgs[arr[1]] = {arr[0]=>{arr[3]=>[arr[4]]}}
						else
							if pkgs[arr[1]].has_key?(arr[0])
								if pkgs[arr[1]][arr[0]].has_key?(arr[3])
									pkgs[arr[1]][arr[0]][arr[3]] << arr[4]
								else
									pkgs[arr[1]][arr[0]][arr[3]] = [arr[4]]
								end
							else
								pkgs[arr[1]][arr[0]] = {arr[3]=>[arr[4]]}
							end
						end
					end
				end
			end
		rescue ThreadError
		end
	end
end

workers.map(&:join)

# {"libpostproc54"=>{"v"=>{"3.2.2-3.4"=>["i586", "x86_64"], "3.2.2-1.1"=>["i586", "x86_64"]}, "i"=>{"3.2.2-3.3"=>["x86_64"]}}}

oss = []
update = []
pkgs.each do |k,v|
	instdver = v["i"].keys[0]
	pkmver = max(v["v"].keys).to_s
	unless v["v"].keys.include?(instdver)
		# if the installed package's release number is near to packman's
		# then you have outdated packman version. if near to oss's
		# then it is oss package.
		instdrel = instdver.gsub(/^.*-/,'')
		vers = v["v"].keys.map {|i| i.gsub(/^.*-/,'')}
		diff = vers.map{|j| j.to_f - instdrel.to_f }.map do |m| 
			if m < 0
				m = m*(-1)
			else
				m = m
			end
		end
		mini = min(diff)
		pos = diff.find_index(mini)
		near = v["v"].keys[pos]
		if near == pkmver
			update << k
		else
			oss << k
		end
	else
		# not from packman
		if instdver != pkmver
			oss << k
		end
	end
end

puts "======================= Packages not from Packman ========================="

if oss.empty?
	puts "Good! All packages are from Packman!"
else
	str = ""
	oss.each do |i| 
		puts i
		# zypper install libavformat57-3.2.2-3.4.x86_64
		pkm = max(pkgs[i]["v"].keys).to_s
		arch = pkgs[i]["i"].values[0][0]
		full = i + "-" + pkm + "." + arch
		str = str + " " + full
	end
	puts "\nFIX: Run 'sudo zypper install" + str + "'."
end

puts "======================= Packman Packages need updates ====================="

if update.empty?
	puts "Good! All packages are from Packman and at their latest versions!"
else
	str = ""
	update.each do |i| 
		puts i
		str = str + " " + i
	end
	puts "\nFIX: Run 'sudo zypper up" + str + "'."
end
