import matplotlib.pyplot as plt
import numpy as np
x = np.linspace(-40, 40, 10000)
y = x ** 2 - 2 * x - 3
plt.figure(1)
plt.title('График функции |y| = x^2 - 2x - 3')
plt.ylabel('Ось y')
plt.xlabel('Ось x')
plt.grid()
plt.axis([-10, 16, 0, 10]) 
plt.plot(x, np.abs(y))
plt.show()