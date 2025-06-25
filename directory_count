python -c "
import os
import json
from collections import defaultdict

def count_files_by_extension(root_dir):
    stats = defaultdict(int)
    directories = 0
    total_files = 0
    
    for root, dirs, files in os.walk(root_dir):
        # Skip hidden directories and common irrelevant dirs
        dirs[:] = [d for d in dirs if not d.startswith('.') and d not in ['__pycache__', 'node_modules', 'build', 'dist']]
        
        directories += len(dirs)
        
        for file in files:
            if not file.startswith('.'):
                ext = os.path.splitext(file)[1].lower()
                if not ext:
                    ext = 'no_extension'
                stats[ext] += 1
                total_files += 1
    
    return dict(stats), directories, total_files

# Count files in directory
stats, dirs, total = count_files_by_extension('.')
print('== Codebase Statistics ===')
print(f'Total Directories: {dirs}')
print(f'Total Files: {total}')
print()
print('Files by Extension:')
for ext, count in sorted(stats.items(), key=lambda x: x[1], reverse=True):
    print(f'  {ext:15} : {count:4d}')
"
