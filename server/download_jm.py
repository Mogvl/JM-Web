#!/usr/bin/env python3
"""
JMComic 图片下载脚本
使用官方 jmcomic 库处理图片重组和格式转换
支持全部章节 / 指定章节下载
"""
import sys, os, json, logging

# 禁止 jmcomic 库的日志输出
logging.disable(logging.CRITICAL)
os.environ['JM_LOG_LEVEL'] = 'ERROR'


def download(jm_id, output_dir, image_format='jpg', chapter_ids=None):
    """使用 jmcomic 库下载并正确处理图片"""
    from jmcomic import JmOption
    from jmcomic import new_downloader

    # 确保使用绝对路径
    output_dir = os.path.abspath(output_dir)
    os.makedirs(output_dir, exist_ok=True)

    # 配置下载选项 - 始终下载 webp（自动处理 descramble）
    option = JmOption.default()
    option.dir_rule.base_dir = output_dir
    option.dir_rule.rule_dsl = 'Bd_Aid'
    option.image_format = 'webp'
    # 关闭多余日志（部分版本仍会打印）
    try:
        from jmcomic import JmModuleConfig
        JmModuleConfig.set_log_file()
    except Exception:
        pass

    if not chapter_ids:
        # 全部章节：直接官方 download_album（最稳）
        from jmcomic import download_album
        download_album(jm_id, option)
    else:
        # 指定章节：拿 album 后按章节构造 photo 逐个下载
        ids = [str(x) for x in chapter_ids]
        client = option.new_jm_client()
        album = client.get_album_detail(jm_id)
        hit = 0
        with new_downloader(option) as dler:
            for idx, ep in enumerate(album.episode_list):
                # ep = (photo_id, index, title)
                if str(ep[0]) in ids or str(ep[1]) in ids:
                    photo = album.create_photo_detail(idx)
                    dler.download_by_photo_detail(photo)
                    hit += 1
            dler.raise_if_has_exception()
        if hit == 0:
            raise RuntimeError('未匹配到任何章节: ' + ','.join(ids))

    # 如果需要 jpg/png，下载后转换
    if image_format != 'webp':
        from PIL import Image
        for root, dirs, files in os.walk(output_dir):
            for f in files:
                if f.endswith('.webp'):
                    webp_path = os.path.join(root, f)
                    new_name = f.replace('.webp', f'.{image_format}')
                    new_path = os.path.join(root, new_name)
                    try:
                        img = Image.open(webp_path)
                        if image_format in ('jpg', 'jpeg'):
                            img = img.convert('RGB')
                            img.save(new_path, 'JPEG', quality=95)
                        elif image_format == 'png':
                            img.save(new_path, 'PNG')
                        os.remove(webp_path)
                    except Exception:
                        pass  # 转换失败保留 webp

    # 统计结果
    total = 0
    for root, dirs, files in os.walk(output_dir):
        total += len([f for f in files if f.endswith(f'.{image_format}')])

    return total


if __name__ == '__main__':
    if len(sys.argv) < 3:
        print(json.dumps({'error': 'usage: script.py <jm_id> <output_dir> [format] [chapter...]'}))
        sys.exit(1)

    jm_id = sys.argv[1]
    output_dir = sys.argv[2]
    image_format = sys.argv[3] if len(sys.argv) > 3 else 'jpg'
    chapter_ids = sys.argv[4:] if len(sys.argv) > 4 else None

    # 重定向 stderr，只输出 JSON 到 stdout
    old_stderr = sys.stderr

    try:
        count = download(jm_id, output_dir, image_format, chapter_ids)
        sys.stdout.write(json.dumps({'success': True, 'count': count, 'dir': output_dir}))
    except Exception as e:
        sys.stderr = old_stderr
        sys.stdout.write(json.dumps({'success': False, 'error': str(e)}))
        sys.exit(1)