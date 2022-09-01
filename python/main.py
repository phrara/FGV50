# 该代码通过接受指纹识别后的txt文件，提取出其中的指纹信息，并利用其阿里云漏洞库爬取相应的潜在漏洞信息，并最终反馈到前端中
# 输入：含指纹信息的txt文件
# 输出： 含潜在漏洞信息的txt文件
import asyncio
import json
import re
import time
from asyncio.proactor_events import _ProactorBasePipeTransport
from functools import wraps

import aiohttp


# 重写父类，避免出现Event loop is closed报错
def silence_event_loop_closed(func):
    @wraps(func)
    def wrapper(self, *args, **kwargs):
        try:
            return func(self, *args, **kwargs)
        except RuntimeError as e:
            if str(e) != 'Event loop is closed':
                raise

    return wrapper


# 同上，避免报错
_ProactorBasePipeTransport.__del__ = silence_event_loop_closed(_ProactorBasePipeTransport.__del__)

# 请求头
headers = {
    'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4482.0 Safari/537.36 Edg/92.0.874.0',
    'Referer': 'https://avd.aliyun.com/'}


class aliyun:
    def __init__(self):
        # 用来记录查询到的潜在漏洞总数
        self.num = 0

    # 异步协程进行url访问
    async def get_url(self, url):
        async with aiohttp.ClientSession() as client:
            async with client.get(url, headers=headers) as resp:
                if resp.status == 200:
                    return await resp.text()

    # 利用正则匹配获取查询后阿里云漏洞库第一页的潜在漏洞有关数据
    def html_parse(self, html):
        # 正则匹配过程
        match = re.compile(
            '<tr>.*?target="_blank">(.*?)</a></td>.*?<td>(.*?)</td>.*?<button.*?>(.*?)</button>.*?nowrap="nowrap">('
            '.*?)</td>' + '.*?<button.*?title="(.*?)">.*?</button>.*?<button.*?title="(.*?)>.*?</tr>', re.S)
        contents = re.findall(match, html)
        # print(contents)
        # 将获取的潜在漏洞信息按照内容类别划分
        for content in contents:
            yield {
                'cve_id': content[0].strip(),
                'vul_name': content[1],
                'vul_type': content[2].strip(),
                'cve_date': content[3].strip(),
                'is_cve': content[-2].strip(),
                'is_poc': content[-1].strip()
            }
        # content = list(content)
        # return contents

    # # 保存内容到指定的txt文件中
    # def save_content_in_text(self, content):
    #     with open('ali_cve.txt', 'a+') as f:
    #         f.write(content + '\n')
    #         # 潜在漏洞数加一
    #         self.num += 1

    async def main(self, ):
        print("读取服务")
        output = []
        keywords = []
        with open("./json/res.json", mode="r", encoding='utf-8') as f:
            res = json.load(f)
            for item in res:
                keywords.append(item["protocol"])
                if item["idstring"] != "":
                    for server in item["idstring"].split(","):
                        keywords.append(server)
        print('*' * 74)
        # 记录开始查询起始时间
        start_time = time.time()
        print('开始查询潜在漏洞>>>')
        # 对每个指纹特征的关键词进行查询
        print(keywords)
        for keyword in keywords:
            keywordlist = []
            url = f'https://avd.aliyun.com/search?q={keyword}'
            html = await self.get_url(url)

            # 对每个潜在漏洞信息进行依次添加到keyword的列表中
            for content in self.html_parse(html):
                keywordlist.append(content)
                self.num += 1
            output.append({keyword: keywordlist})
        with open("./json/ali_cve.json", "w", encoding='utf-8') as f:
            f.write(json.dumps(output, ensure_ascii=False))
        print(output)
        # 记录最后时间，用来计算查询总时长
        end_time = time.time()
        print(f'本次共查询{self.num}个潜在漏洞，共耗时{end_time - start_time}秒。')


if __name__ == '__main__':
    # 保证ali_cve.txt文件在爬虫获取信息时为空白的
    run_data = aliyun()
    asyncio.run(run_data.main())
