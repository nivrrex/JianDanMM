# -*- coding: utf-8 -*-
import requests
from pyquery import PyQuery as pyq
from lxml import etree
import re
import json
import time

def get_file_name(url):
    pattern = re.compile(r'http.*//(.*)')
    result = pattern.findall(url)
    if result != None:
        return result[0]
    else:
        return None

def get_file_name2(url):
    result = url.replace("https://","")
    result = result.replace("http://","")
    result = result.replace("/","-")
    result = result.replace("?","")
    if result != None:
        return result
    else:
        return None

def get_scale(scale):
    if scale >= 10:
        level = "very"
    elif scale >= 5:
        level = "high"
    elif scale >= 1.5:
        level = "good"
    elif scale >= 0.75:
        level = "normal"
    elif scale < 0.75:
        level = "bad"
    return level


def get_jpg(id, url, support, unsupport):
    scale = float(support) / float(unsupport)
    rank = get_scale(scale)

    #获取具体图片并写入文件
    pic_body = ""
    try:
        r = requests.get(url)
        pic_body = r.content
    except:
        print("******{0} is download error.".format(url))
        return
    else:
        pass
    finally:
        pass

    filename = get_file_name2(url)
    filepath = ("./{0}---{1}").format(rank,filename)

    try:
        f = open(filepath,"wb")
        f.write(pic_body)
    except:
        print("++++++{0} is write error.".format(filename))
        raise
        return
    else:
        pass
    finally:
        f.close()

    print("Write OK ... {0}   {1}\t{2}\t{3}".format(id,filename,support,unsupport))


def get_jiandan_mm_pic(page_num):
    url = 'http://jandan.net/ooxx/page-' + str(page_num)
    html = pyq(url)
    print('reading ...  http://jandan.net/ooxx/page-{0}\n'.format(page_num))
    #print(html)

    hash_pic_message = {}
    #获取图片地址
    for element in html('li div div.row div.text'):
        img = pyq(element).find('img')
        #img = pyq(element)('img')
        if img != None:
            id = pyq(element)('span a').text()
            id = id.replace("vote-","")
            hash_pic_message[id]={}
            hash_pic_message[id]['ID']=id
            hash_pic_message[id]['URL']=[]
            hash_pic_message[id]['FileName']=[]

            if img.attr('org_src') == None:
                for t in img:
                    url = img(t).attr('src')
                    hash_pic_message[id]['URL'].append(url)
                    hash_pic_message[id]['FileName'].append(get_file_name2(url))
            else:
                for t in img:
                    url = img(t).attr('org_src')
                    hash_pic_message[id]['URL'].append(url)
                    hash_pic_message[id]['FileName'].append(get_file_name2(url))
    
    #获取图片ID和评级
    for element in html('li div div.row div.vote'):
        id = pyq(element).attr('id')
        id = id.replace("vote-","")
        
        vote = pyq(element).text()

        reg_vote = 'OO \[ (\d.*) \] XX \[ (\d.*) \]'
        pattern = re.compile(reg_vote)
        result = pattern.findall(vote)
        if result != None:
            support = result[0][0]
            unsupport = result[0][1]
            hash_pic_message[id]["Support"] = support
            hash_pic_message[id]["UnSupport"] = unsupport
            
            scale =  float(support) / float(unsupport)
            rank = get_scale(scale)
            hash_pic_message[id]["Scale"] = scale
            hash_pic_message[id]["Rank"] = rank
            
    for value in hash_pic_message.values():
        #print(value)
        pass
    return hash_pic_message.values()

if __name__=="__main__":
    #reading config.json file
    f = open('config.json','r')
    data = json.load(f)
    startpage = data['startpage']
    endpage = data['endpage']

    for page in range(int(startpage),int(endpage)):
        for url_index in get_jiandan_mm_pic(page):
            #print (url_index,"\n")
#            print(url_index["URL"][0])
            for url in url_index["URL"]:
                get_jpg(url_index["ID"],url,url_index["Support"],url_index["UnSupport"])
            time.sleep(0.2)
        print("\n")
        time.sleep(0.8)




#################################OLD#################################
#def get_jiandan_mm_pic(page_num):
#    url = 'http://jandan.net/ooxx/page-' + str(page_num)
#    r = requests.get(url)
#    pic_body = r.text
#    return_status_code = r.status_code
#    print('reading ...  http://jandan.net/ooxx/page-{0}   :   Return Code -- {1} . \n'.format(page_num,return_status_code))
#    #print(pic_body.encode("utf-8"))
#
#    hash_pic_message = {}
#
#    reg_str_jpg = 'comment-(\d*?)">\d*?</a></span><p>.*?\s*<img src="([a-zA-z]+://[^\s]+?)" /></p>'
#    pattern_jpg = re.compile(reg_str_jpg)
#    result_jpg = pattern_jpg.findall(pic_body)
#    for tdata_jpg in result_jpg:
#        if not tdata_jpg[0] in hash_pic_message:
#            hash_pic_message[tdata_jpg[0]]={}
#            hash_pic_message[tdata_jpg[0]]['URL']=tdata_jpg[1]
#            hash_pic_message[tdata_jpg[0]]['FileName']=get_file_name2(tdata_jpg[1])
#            #print(tdata_jpg)
#        else:
#            print("...error... already has this picture. {0}".format(tdata_jpg))
#
#    reg_str_gif = 'comment-(\d*?)">\d*?</a></span><p>.*?\s*<img src="[a-zA-z]+://[^\s]+?" org_src="([a-zA-z]+://[^\s]+?)" o'
#    pattern_gif = re.compile(reg_str_gif)
#    result_gif = pattern_gif.findall(pic_body)
#    for tdata_gif in result_gif:
#        if not tdata_gif[0] in hash_pic_message:
#            hash_pic_message[tdata_gif[0]]={}
#            hash_pic_message[tdata_gif[0]]['URL']=tdata_gif[1]
#            hash_pic_message[tdata_gif[0]]['FileName']=get_file_name2(tdata_gif[1])
#            #print(tdata_gif)
#        else:
#            print("...error... already has this picture. {0}".format(tdata_gif))
#
#    reg_str_vote = '<div class="vote" id="vote-(\d*?)">.*?<span id="cos_support-\d*?">(\d*?)<\/span>.*?<span id="cos_unsupport-\d*?">(\d*?)<\/span>'
#    pattern_vote = re.compile(reg_str_vote)
#    result_vote = pattern_vote.findall(pic_body)
#    for tdata_vote in result_vote:
#        if tdata_vote[0] in hash_pic_message:
#            support = tdata_vote[1]
#            unsupport = tdata_vote[2]
#            scale =  float(support) / float(unsupport)
#            rank = get_scale(scale)
#            hash_pic_message[tdata_vote[0]]["Support"] = support
#            hash_pic_message[tdata_vote[0]]["UnSupport"] = unsupport
#            hash_pic_message[tdata_vote[0]]["Scale"] = scale
#            hash_pic_message[tdata_vote[0]]["Rank"] = rank
#        else:
#            print("...error... no already has this picture. {0}".format(tdata_vote))
#
#    for value in hash_pic_message.values():
#        #print(value)
#        pass
#    return hash_pic_message.values()
##################################################################