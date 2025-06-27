python -c "
import os
import glob

# Find all major component directories
components = {
    'Frontend': ['frontend/', 'platform/src/components/', 'frontend/src/'],
    'Backend': ['core/backend/', 'platform/modules/backend/', 'apps/api/'],
    'Python Modules': ['platform/modules/', 'platform/src/', 'core/', 'libs/'],
    'Rust Components': ['platform/rust/', 'platform/modules/agent/', 'platform/modules/defensive_security/'],
    'Go Components': ['platform/modules/net/', 'platform/modules/cloud/'],
    'Tests': ['tests/', 'frontend/src/**/__tests__/', 'platform/modules/backend/__tests__/'],
    'Documentation': ['docs/', 'README.md'],
    'Configuration': ['configs/', '.env*', '*.toml', '*.json', 'docker-compose.yml'],
    'Infrastructure': ['infrastructure/', 'deployments/', '.github/'],
    'Scripts': ['scripts/', 'tools/scripts/'],
    'CLI': ['cli/'],
    'Database': ['migrations/', '*.db', 'core/database/'],
}

print('=== Platform Component Index ===')
for component, paths in components.items():
    print(f'\n{component}:')
    found_paths = []
    for pattern in paths:
      matches = glob.glob(pattern, recursive=True)
        if matches:
            found_paths.extend([m for m in matches if os.path.exists(m)])
    
    if found_paths:
        for path in sorted(set(found_paths)):
            if os.path.isdir(path):
                file_count = sum(len(files) for _, _, files in os.walk(path))
                print(f'  {path} ({file_count} files)')
            elif os.path.isfile(path):
                print(f'  {path}')
    else:
        print('  No components found')
"
