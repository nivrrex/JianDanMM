#encoding:gbk
require 'net/http'

def get_file_name(url)
    re = /http.*\/(.*)/.match(url)
    if re != nil
        return re[1]
    else
        return nil
    end
end

def get_file_name2(url)
    re = url.gsub("https://","")
    re = re.gsub("http://","")
    re = re.gsub("/","-")
    re = re.gsub("?","")
    if re != nil
        return re
    else
        return nil
    end
end

def get_scale(scale)
    if scale >= 10 then
        level = "very"
    elsif scale >= 5 then
        level = "high"
    elsif scale >= 1.5 then
        level = "good"
    elsif scale >= 0.75 then
        level = "normal"
    elsif scale < 0.75 then
        level = "bad"
    end
    return level
end


def get_jpg(url , file_name ,id , support , unsupport ,rank)
    log = ""
    begin
        redirect = 1
        redirecting = false 
        begin
            uri = URI.parse(url) 
            req = Net::HTTP::Get.new(uri.path) 
            res = Net::HTTP.start(uri.host, uri.port) do |http| 
                http.read_timeout = 10
                http.request(req) 
            end 
            if res.header['location'] # 遇到重定向，则url设定为location，再次抓取 
                url = res.header['location']  
                redirecting = true 
            end 
            redirect -= 1 
        end while redirecting and redirect>=0 
    rescue Timeout::Error
        log = "****** Http Error: #{id} #{url} is not get."
        return log
    rescue Errno::ETIMEDOUT
        log = "****** Http Error: #{id} #{url} is not open."
        return log
    rescue Errno::ECONNRESET
        log = "****** Http Error: #{id} #{url} is 'An existing connection was forcibly closed by the remote host.'."
        return log
    rescue Errno::EADDRNOTAVAIL
        log = "****** Http Error: #{id} #{url} is 'The requested address is not valid in its context.'."
        return log
    rescue ArgumentError
        log = "****** Http Error: #{id} #{url} is 'HTTP request path is empty'."
        return log
    rescue SocketError
        log = "****** Http Error: #{id} #{url} is 不知道这样的主机。."
        return log
    rescue EOFError
        log = "****** Http Error: #{id} #{url} is end of file reached."
        return log
    rescue Errno::ECONNREFUSED
        log = "****** Http Error: #{id} #{url} is 'No connection could be made because the target machine actively refused it'."
        return log
    rescue URI::InvalidURIError
        log = "****** Http Error: #{id} #{url} is bad URI."
        return log
    rescue
        log = "****** Http Error: #{id} #{url} is other error."
        return log
    end
  
    begin
        open(file_name, "wb") do |file|
            file.write(res.body)
        end
    rescue
        log = "++++++ Write Error: #{id} #{url} is write error."
        return log
    end
 
    log = "Write OK ... #{id} #{url} #{file_name} #{support} #{unsupport} #{rank}"
    
    return log
end

$url_hash = {}
f_to = File::open("support.result.log","w")

900.upto (901) do |i|
    uri = URI.parse('http://jandan.net/ooxx/page-' + i.to_s())  
    http = Net::HTTP.new(uri.host, uri.port)
    request = Net::HTTP::Get.new(uri.request_uri)
    request.initialize_http_header({"User-Agent" => "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Maxthon/4.4.6.2000 Chrome/30.0.1599.101 Safari/537.37#{i}"})

    response = http.request(request)
    head = response.code
    html = response.body    

    puts 
    log_str = "reading " + 'http://jandan.net/ooxx/page-' + i.to_s() + "   :   Return Code -- " + head
    puts log_str
	f_to.puts(log_str)
    puts
                
    html.scan(/comment-(\d*?)">\d*?<\/a><\/span><p>.*?\s*<img src="([a-zA-z]+:\/\/[^\s]+?)" \/><\/p>/m) do |data|
        id = data[0]
        url = data[1]
        file_name = get_file_name2(url)
        if not $url_hash.has_key?(id) then
            $url_hash[id] = {}
        else
            #handle error
        end
        $url_hash[id]["URL"] = url
        $url_hash[id]["FileName"] = file_name
    end

    html.scan(/comment-(\d*?)">\d*?<\/a><\/span><p>.*?\s*<img src="[a-zA-z]+:\/\/[^\s]+?" org_src="([a-zA-z]+:\/\/[^\s]+?)" o/m) do |data|
        id = data[0]
        url = data[1]
        file_name = get_file_name2(url)
        if not $url_hash.has_key?(id) then
            $url_hash[id] = {}
        else
            #handle error
        end
        $url_hash[id]["URL"] = url
        $url_hash[id]["FileName"] = file_name
    end

    html.scan(/<div class="vote" id="vote-(\d*?)">.*?<span id="cos_support-\d*?">(\d*?)<\/span>.*?<span id="cos_unsupport-\d*?">(\d*?)<\/span>/m) do |data|
        id = data[0]
		support = data[1]
		unsupport = data[2]
		scale = support.to_f / unsupport.to_f
		rank = get_scale(scale)

        if not $url_hash.has_key?(id) then
            $url_hash[id] = {}
        else
            log_str = "%%%%%%  HTML regexp error ... #{id}  #{$url_hash[id]["URL"]}has no .."
            #handle error
        end
		$url_hash[id]["Support"] = support
		$url_hash[id]["UnSupport"] = unsupport
		$url_hash[id]["Scale"]  = scale.to_s
		$url_hash[id]["Rank"]  = rank
        if $url_hash[id]["FileName"] != nil
            $url_hash[id]["FileName"]  = rank + "---" + $url_hash[id]["FileName"]
            log_str = get_jpg($url_hash[id]["URL"] ,$url_hash[id]["FileName"], id ,support ,unsupport ,rank)
        else
            log_str = "%%%%%%  HTML regexp error ... #{id}  is error"
        end
		puts log_str
		f_to.puts(log_str)

        STDOUT.flush
        sleep Random.rand(0.01..0.05)
    end
    STDOUT.flush
    sleep Random.rand(0.8..1.0)
end
f_to.close
