#!/usr/bin/env python3
"""
JMComic 图片下载脚本
使用官方 jmcomic 库处理图片重组和格式转换
"""
import sys, os, json

def download(jm_id, output_dir, image_format='jpg'):
    """使用 jmcomic 库下载并正确处理图片"""
    from jmcomic import JmOption, download_album, JmModuleConfig

    # 配置下载选项
    option = JmOption.default()
    option.dir_rule.download_dir = output_dir
    option.dir_rule.base_rule = 'Bd_Aid'
    option.image_format = image_format

    # 下载（自动处理 descramble）
    download_album(jm_id, option)

    # 统计结果
    total = 0
    for root, dirs, files in os.walk(output_dir):
        total += len([f for f in files if f.endswith(f'.{image_format}')])

    return total

if __name__ == '__main__':
    if len(sys.argv) < 3:
        print(json.dumps({'error': 'usage: script.py <jm_id> <output_dir> [format]'}))
        sys.exit(1)

    jm_id = sys.argv[1]
    output_dir = sys.argv[2]
    image_format = sys.argv[3] if len(sys.argv) > 3 else 'jpg'

    try:
        count = download(jm_id, output_dir, image_format)
        print(json.dumps({'success': True, 'count': count, 'dir': output_dir}))
    except Exception as e:
        print(json.dumps({'success': False, 'error': str(e)}))
        sys.exit(1)