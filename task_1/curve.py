import numpy as np
import matplotlib.pyplot as plt

f2 = open('group.txt', 'r')
l = 2
ll = 3
lline= f2.readlines()
lline_coor_x= lline[l].split(",")
lline_coor_y = lline[ll].split(",")

print(lline_coor_x)
print(lline_coor_y)
f2.close()

result_x = [int(item) for item in lline_coor_x[:-1]]
result_y = [int(item) for item in lline_coor_y[:-1]]


plt.figure(figsize=(8, 6))
_, ax = plt.subplots()
ax.scatter(result_x, result_y)
plt.show()

