rm -rf bullet/Bullet-C-Api.h
rm -rf bullet/BulletMultiThreaded
rm -rf bullet/MiniCL
rm -rf bullet/BulletCollision/Gimpact
rm -rf bullet/vectormath

find ./bullet ! -name "*.h" -type f -delete
find ./bullet -type d -name CMakeFiles -exec rm -rf {} \;

