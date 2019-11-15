import matplotlib.pyplot as plt
from mpl_toolkits.mplot3d import Axes3D


# fig = plt.figure()
# ax = fig.add_subplot(111, projection='3d')

with open('search_space.txt') as f:
    points = f.readlines()

max_dist = float('-inf')
xs, ys, zs = [], [], []
for point in points:
    x, y, z = map(float, point.split())
    d = (x**2+y**2+z**2)**.5
    if d > max_dist:
        print(f'new max dist: {d}>{max_dist}')
        max_dist = d
    xs.append(x)
    ys.append(y)
    zs.append(z)
print(max(xs))
print(max(ys))
print(max(zs))
# ax.plot(xs, ys, zs, 'bo', markersize=.01, linestyle='none')
# plt.show()
