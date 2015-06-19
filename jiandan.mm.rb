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

def get_jpg(url)
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
        puts "****** #{url} is not get."
        return
    rescue Errno::ETIMEDOUT
        puts "****** #{url} is not open."
        return
    rescue Errno::ECONNRESET
        puts "****** #{url} is 'An existing connection was forcibly closed by the remote host.'."
        return
    rescue Errno::EADDRNOTAVAIL
        puts "****** is 'The requested address is not valid in its context.'."
        return
    rescue ArgumentError
        puts "****** #{url} is 'HTTP request path is empty'."
        return
    rescue SocketError
        puts "****** #{url} is 不知道这样的主机。."
        return
    rescue EOFError
        puts "****** is end of file reached."
        return
    rescue Errno::ECONNREFUSED
        puts "****** is 'No connection could be made because the target machine actively refused it'."
        return
    rescue URI::InvalidURIError
        puts "****** is bad URI."
        return
    rescue
        puts "****** is other error."
        return
    end
 
    filename = get_file_name2(url)
 
    begin
        open("./test/" + filename, "wb") do |file|
            file.write(res.body)
        end
    rescue
        puts "++++++ is write error."
        return
    end
 
    puts "Write OK ... #{filename}"
end


f_to = File::open("support.result.csv","w")
900.upto (1430) do |i|
    uri = URI.parse('http://jandan.net/ooxx/page-' + i.to_s())  
    http = Net::HTTP.new(uri.host, uri.port)
    request = Net::HTTP::Get.new(uri.request_uri)
    request.initialize_http_header({"User-Agent" => "Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; rv:11.0) like Gecko"})

    response = http.request(request)
    head = response.code
    html = response.body    

    puts
    puts "reading " + 'http://jandan.net/ooxx/page-' + i.to_s() + "   :   Return Code -- " + head
    puts

    html.scan(/<li id="comment-\d*?">.*?<img src="(.*?)".*?<span id="cos_support-\d*?">(\d*?)<\/span>.*?<span id="cos_unsupport-\d*?">(.*?)<\/span>.*?<\/li>/m) do |data|
        result = data[0] + " " + get_file_name2(url) + " " + data[1] + " " + data[2]
        puts result
        f_to.puts(result)
        get_jpg (data[0])
        sleep Random.rand(0.01..0.05)
        STDOUT.flush
    end
    sleep Random.rand(0.8..1.0)
    STDOUT.flush
end
f_to.close
