#!/usr/bin/env ruby
#
# Best Mirror: find best openSUSE/Packman mirror for you
#
# usage: ruby bstmirror.rb -region="Asia" -os=422 -file=quick
# 	-region: Africa, Asia, Europe, North America, South America, Oceania
# 	-os: tw, 422, 421
# 	-file: quick, long

require 'nokogiri'
require 'open-uri'
require 'net/http'
require 'uri'
require 'timeout'
require 'fileutils'

def resp(uri="")
	Timeout::timeout(5) do
		res = Net::HTTP.get_response(URI(uri))
		return true
	end
rescue Timeout::Error
	return false
end

def speedtest(url="")
	# http://mirrors.hust.edu.cn/packman/suse/openSUSE_Tumbleweed/Essentials/repodata/primary.xml.gz
	site = url.gsub(/^http:\/\//,'').gsub(/\/.*$/,'')
	file_name = site + "_" + url.gsub(/^.*\//,'')
	path = url.gsub(/^http:\/\//,'').gsub(site,'')
	st = Time.now
	f = open(file_name,'w')
	Net::HTTP.start(site) do |http|
		begin
			http.request_get(path) do |resp|
				resp.read_body {|segment| f.write(segment)}
			end
		ensure
			f.close
		end
	end
	et = Time.now
	duration = et - st
	size = File.size(file_name) / 1024.0
	speed = size / duration
	FileUtils.rm_rf file_name
	# KB/s
	return speed
end

def max(arr=[])
	size = arr.length
	max = arr[size-1]
	size.times do |i|
		if arr[i] > max
			max = arr[i]
		end
	end
	return max
end

def optParser()
	opts = {}
	ARGV.each do |arg|
		arg_arr = arg.split("=")
		opts[arg_arr[0].gsub("-","")] = arg_arr[1]
	end
	return opts
end

oss = Nokogiri::HTML(open("http://mirrors.opensuse.org"))
pm = Nokogiri::HTML(open("http://packman.links2linux.de/mirrors"))
oss_world,pm_world = {},[]
continent = ""
opts = optParser
checkfile = "repomd.xml"

# data structure
#
# {"Asia"=>{"China"=>{"Sohu"=>{"HTTP"=>"","Priority"=>3,"TW"=>true,"422"=>true,"421"=>true}},"Taiwan"=>{}}}

oss.xpath("//table/tbody/tr").each do |tr|
	# <tr>
	#   <td colspan="31" class="newregion">Africa:</td>
	# </tr>
	unless tr.xpath("td").attribute("class").nil?
		continent = tr.xpath("td").text.gsub(":","")
	else
		# 0 <td>Ecuador</td>
		# 1 <td><a href="http://www.cedia.org.ec">Consorcio Ecuatoriano para el Desarrollo de Internet Avanzado</a></td>
		# 2 <td><a href="http://mirror.cedia.org.ec/opensuse">HTTP</a></td>
		# 3 <td><a href="ftp://mirror.cedia.org.ec/opensuse/">FTP</a></td>
		# 4 <td><a href="rsync://mirror.cedia.org.ec/opensuse/">rsync</a></td>
		# 5 <td>***</td>
		# 6 <td class="a"></td>
		# 7 <td class="b"></td> TW
		# 8 <td class="a">√</td>
		# 9 <td class="b">√</td> 422
		# 10 <td class="a"></td>
		# 11 <td class="b"></td>
		# 12 <td class="a">√</td> 422 update 
		# 13 <td class="b">√</td>
		# 14 <td class="a">√</td> 421
		# 15 <td class="b"></td>
		# 16 <td class="a"></td>
		# 17 <td class="b">√</td> 421 update
		td0 = tr.at_xpath("td[1]").text
		if td0.length > 0 # ignore the placeholders
			td1,td2,td3 = tr.at_xpath("td[2]/a").text,nil,nil
			td2 = tr.at_xpath("td[3]/a/@href").value unless tr.at_xpath("td[3]/a").nil?
			td3 = tr.at_xpath("td[4]/a/@href").value unless tr.at_xpath("td[4]/a").nil? 
			td5 = tr.at_xpath("td[6]").text.length
			td7 = tr.at_xpath("td[8]").text.length
			td9 = tr.at_xpath("td[10]").text.length
			td14 = tr.at_xpath("td[15]").text.length

			country,name,http,ftp,pri,tw,o422,o421 = tr.at_xpath("td[1]").text.strip!,nil,td2,td3,td5,true,true,true
			if td1.length != 0
				name = td1
			else
				name = td2
			end
			tw = false if td7 == 0
			o422 = false if td9 == 0
			o421 = false if td14 == 0
			
			if oss_world.has_key? continent
				if oss_world[continent].has_key? country
					oss_world[continent][country][name] = {"HTTP"=>http,"FTP"=>ftp,"pri"=>pri,"tw"=>tw,"o422"=>o422,"o421"=>o421}
				else
					oss_world[continent][country] = {name=>{"HTTP"=>http,"FTP"=>ftp,"pri"=>pri,"tw"=>tw,"o422"=>o422,"o421"=>o421}}
				end
			else
				oss_world[continent] = {country=>{name=>{"HTTP"=>http,"FTP"=>ftp,"pri"=>pri,"tw"=>tw,"o422"=>o422,"o421"=>o421}}}
			end
		end
	end
end

oss_avail_mirrors,oss_resp_mirrors,flavor = [],[],""
if opts.has_key?("os")
	if opts["os"] == "422"
		flavor = "o422"
	elsif opts["os"] == "421"
		flavor = "o421"
	else
		flavor = "tw"
	end
else
	flavor = "o422"
end

if opts.has_key?("region")
	oss_world[opts["region"]].each_value do |v|
		v.each_value do |w|
			if w[flavor]
				unless w["HTTP"].nil?
					oss_avail_mirrors << w["HTTP"]
				else
					oss_avail_mirrors << w["FTP"]
				end
			end
	        end	       
	end
else
	oss_world.each_value do |v|
		v.each_value do |w|
			w.each_value do |x|
				if x[flavor]
					unless x["HTTP"].nil?
						oss_avail_mirrors << x["HTTP"]
					else
						oss_avail_mirrors << x["FTP"]
					end
				end
			end
		end
	end
end

oss_avail_mirrors.each {|i| oss_resp_mirrors << i if resp(i)}

oss_mirror_speeds = {}
oss_pool_size = oss_resp_mirrors.length
oss_jobs = Queue.new
oss_resp_mirrors.each {|i| oss_jobs.push i}
oss_workers = oss_pool_size.times.map do
	Thread.new do
		begin
			while x = oss_jobs.pop(true)
				y = x
				y = x + "/" unless y.index(/\/$/)
				checkfile = "appdata.xml.gz" if ( opts.has_key?("file") && opts["file"] == "long" )
				url = y + "tumbleweed/repo/oss/suse/repodata/" + checkfile
				sp = speedtest(url)
				oss_mirror_speeds[x] = sp
			end
		rescue ThreadError
		end
	end
end

oss_workers.map(&:join)

best_oss = oss_mirror_speeds.key(max(oss_mirror_speeds.values))

# Packman

pm.xpath('//td[@class="mirrortable mirror"]').each do |td|
	unless td.at_xpath("a").nil? # ignore rsync mirror
		v = td.at_xpath("a/@href").value
		pm_world << v unless v.index("ftp://")
	end
end

pm_resp_mirrors = []
pm_world.each {|i| pm_resp_mirrors << i if resp(i)}

pm_mirror_speeds = {}
pm_pool_size = pm_resp_mirrors.length
pm_jobs = Queue.new
pm_resp_mirrors.each {|i| pm_jobs.push i}
pm_workers = pm_pool_size.times.map do
	Thread.new do
		begin
			while x = pm_jobs.pop(true)
				y = x
				y = x + "/" unless y.index(/\/$/)
				checkfile = "primary.xml.gz" if ( opts.has_key?("file") && opts["file"] == "long")
				url = y + "suse/openSUSE_Tumbleweed/Essentials/repodata/" + checkfile
				sp = speedtest(url)
				pm_mirror_speeds[x] = sp
			end
		rescue ThreadError
		end
	end
end

pm_workers.map(&:join)

best_pm = pm_mirror_speeds.key(max(pm_mirror_speeds.values))

puts "======================= Best openSUSE Mirror ======================"
puts best_oss
puts "======================= Best Packman Mirror ======================="
puts best_pm




