import os
import datetime

HEADER_TEMPLATE = """/*
 * @Author: JimZhang
 * @Date: {date}
 * @LastEditors: JimZhang
 * @LastEditTime: {time}
 * @FilePath: /{filepath}
 * @Description: 
 * 
 */
"""

def add_headers_to_files(root_dir):
    extensions = {'.go', '.vue', '.ts', '.js'}
    for subdir, dirs, files in os.walk(root_dir):
        if 'node_modules' in subdir or '.git' in subdir or 'dist' in subdir:
            continue
        for file in files:
            ext = os.path.splitext(file)[1]
            if ext in extensions:
                filepath = os.path.join(subdir, file)
                try:
                    with open(filepath, 'r+', encoding='utf-8') as f:
                        content = f.read()
                        if not content.startswith('/*\n * @Author:'):
                            # Generate header
                            rel_path = os.path.relpath(filepath, root_dir)
                            now = datetime.datetime.now()
                            date_str = now.strftime('%Y-%m-%d %H:%M:%S')
                            header = HEADER_TEMPLATE.format(
                                date=date_str,
                                time=date_str,
                                filepath=rel_path.replace('\\', '/')
                            )
                            f.seek(0, 0)
                            f.write(header + content)
                            print(f'Added header to {rel_path}')
                except Exception as e:
                    print(f'Failed {filepath}: {e}')

if __name__ == "__main__":
    add_headers_to_files('/Users/jim/Desktop/Code/Project/general_go_vue_admin')
