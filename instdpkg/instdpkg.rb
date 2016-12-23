#!/usr/bin/env ruby
#
# return installed packages during the user-specified date/time
#
# usage: ./instdpkg -d="2016-12-18" -t="00:01:31"
# 	-tl:	just tell me the days I install/remove packages.
#	-d:	check that date. if omitted, it'll check today by default.
#	-t:	check from that time. if omitted, it'll check from 00:00:00

def readLog()
	pkgs,tl = {},[]
	open("/var/log/zypp/history") do |f|
		f.each_line do |line|
			if line.index(/install\||remove\s\|/) && ! line.index("|_")
				arr = line.split("|")
				# "2016-12-23 00:01:31" "install" "vlc" "0.1.0-1.1" "x86_64" "packman"
				date = arr[0].gsub(/\s.*$/,'')
				time = arr[0].gsub(/^.*\s/,'')
				repo = ""
				if arr[1] != "install"
					arr[1] = "remove"
					repo = "none"
				else
					repo = arr[6]
				end
				if pkgs.has_key?(date)
					#"2016-12-23"=>{"vlc"=>["00:01:31","install","0.1.0-1.1.x86_64","packman"]}
					pkgs[date][arr[2]] = [time,arr[1],arr[3] + "." + arr[4],repo]
				else
					pkgs[date] = {arr[2]=>[time,arr[1],arr[3] + "." + arr[4],repo]}
				end
				tl << date 
			end	
		end
	end

	tl = tl.uniq.sort

	return pkgs,tl
end

def optionParser()
	opts = Hash.new
	ARGV.each do |opt|
		arr = opt.split("=")
		arr[0] = arr[0].gsub("-","")
		opts[arr[0]] = arr[1]
	end
	return opts
end

def comptime(old="",new="")
	oa,na = old.split("\s"),new.split("\s")
	oad,oat = oa[0].split("-"),oa[1].split(":")
	nad,nat = na[0].split("-"),na[1].split(":")
	od = Time.utc(oad[0].to_i,oad[1].to_i,oad[2].to_i,oat[0].to_i,oat[1].to_i,oat[2].to_i)
	nd = Time.utc(nad[0].to_i,nad[1].to_i,nad[2].to_i,nat[0].to_i,nat[1].to_i,nat[2].to_i)
	if nd.to_i >= od.to_i
		return true
	else
		return false
	end
end

pkgs,tl,opts,date,time = readLog[0],readLog[1],optionParser,"",""
date = opts["d"] if opts.has_key?("d")
time = opts["t"] if opts.has_key?("t")
today = Time.now.strftime("%Y-%m-%d %H:%M:%S")
date = today.gsub(/[\s].*$/,'') if date.empty?
time = "00:00:00" if time.empty?
dt = date + " " + time

if opts.has_key?("tl")
	tl.each {|d| puts d}
else
	unless pkgs[date].nil?
		if opts.has_key?("t")
            puts "================ Packages altered on " + date + " after " + time + " ======================"
		    puts "|\s\sdate\s\s\s|\s\stime\s\s|\soper.\s| name | version | repo |"
			pkgs[date].each do |k,v|
				# {"ffmpeg"=>["23:26:00", "install", "3.2.2-3.3.x86_64", "packman"]}
				new_dt = date + "\s" + v[0]
				if comptime(dt,new_dt)
					puts date + "|" + v[0] + "|" + v[1] + "|" + k + "|" + v[2] + "|" + v[3]
				end
			end
        else
            puts "================ Packages altered on " + date + " ======================"
		    puts "|\s\sdate\s\s\s|\s\stime\s\s|\soper.\s| name | version | repo |"
			pkgs[date].each do |k,v|
				puts date + "|" + v[0] + "|" + v[1] + "|" + k + "|" + v[2] + "|" + v[3]
			end
		end
	else
		puts "you didn't install/remove packages today. forgot the \"-d\" option?"
	end
end
