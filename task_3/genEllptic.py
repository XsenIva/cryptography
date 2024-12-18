import sympy
import random
import matplotlib.pyplot as plt

primes = set()

def generate_prime(l):
    if l <= 4:
        print("Cлишком маленькое p")
        return None

    while True:
        p = sympy.randprime(2**(l-1), 2**l)
        if p % 4 == 1 and p not in primes:
            primes.add(p)
            return p

def get_k(p, ai, q):
    k = 0
    while pow(ai, (2 ** k) * q, p) != 1:
        k += 1
    return k


def div_to_squar(p, D):
    quad_residue = sympy.is_nthpow_residue(sympy.sqrt_mod(-D, p), 2, p)
    if quad_residue == False:
        return None

    u = sympy.sqrt_mod(-D, p)
    i, u, m = 0,  [u],  [p]

    while m[-1] != 1:
        m.append((u[i] * u[i] + D) / m[i])
        u.append(min(
            u[i] % m[i + 1],
            m[i + 1] - u[i] % m[i + 1]
        ))
        if (m[-1] == 1):
            break
        i = i + 1

    a, b = u[-2], 1
    while i > 0:

        aFirst  = (u[i - 1] * a + D * b) / (a*a + D*b*b)
        aSecond = (- u[i - 1] * a + D * b) / (a*a + D*b*b)

        bFirst  = (-a + u[i - 1] * b) / (a*a + D*b*b)
        bSecond = (-a - u[i - 1] * b) / (a*a + D*b*b)

        if int(aFirst) == aFirst:
            a = aFirst
        else:
            a = aSecond

        if int(bFirst) == bFirst:
            b = bFirst
        else:
            b = bSecond

        i = i - 1

    return a, b


def is_int(a):
    return int(a) == a

def prov_t_cond(p, a, b, m):
    list_t_arg = [ 2*a, -2*a, 2*b, -2*b 
    ]
    coff_r = [2, 4]
    pair_n_r = []
    result = []

    for t in list_t_arg:
        N = p + 1 + t
        for r_t in coff_r:
            r = N / r_t
            if (is_int(r) and sympy.isprime(int(r))):
                pair_n_r.append((N, r))
    
    for pair in pair_n_r:
        pair_n, pair_r = pair
        prover = True
        for i in range(m):
            prover = prover and p != int(pair_r) and pow(p, i + 1, int(pair_r)) != 1
        if prover:
            result.append((p, pair_n, pair_r))

    return result


def add_ellip_points(x, y, x1, y1, A, p):
    if x == None and y == None:
        return x1, y1

    if x1 == None and y1 == None:
        return x, y
    
    x = x % p
    x1 = x1 % p
    y = y % p
    y1 = y1 % p

    alpha = None
    if x != x1 or y != y1:
        if x == x1:
            return None, None
        alpha = ((y1 - y) * sympy.mod_inverse(x1 - x, p)) % p
    else:
        if y == 0:
            return None, None
        alpha = ((3 * x * x + A) * sympy.mod_inverse(2 * y, p)) % p

    x3 = (alpha ** 2 - x - x1) % p
    y3 = (alpha * (x - x3) - y) % p

    return (x3, y3)


def scalarMultiply(other, x, y, p, A):
    if other == 0:
        return (None, None)
    negate = False
    if other < 0:
        negate = True
        other = -other
    init = False
    powerOf2 = (x, y)
    while True:
        if other % 2 == 1:
            if init:
                res = add_ellip_points(*res, *powerOf2, A, p)
            else:
                res = powerOf2
                init = True
        other //=2
        if other == 0:
            break
        powerOf2 = add_ellip_points(*powerOf2, *powerOf2, A, p)
    if negate:
        res = (res[0], -res[1] % p)
    return res


def check_infinite_point(N, x, y, A, p):
    return scalarMultiply(N, x, y, p, A) == (None, None)


def generate_points(p, N, r, points : set):
    while len(points) != N:
        x = random.randint(1, p - 1) 
        y = random.randint(1, p - 1) 
        if (x, y) in points:
            continue
        points.add((x, y))
       
        a_coff = ((y * y - x * x * x)*sympy.mod_inverse(x, p)) % p 
        
        if not is_int(a_coff): continue

        quad_residue = sympy.is_nthpow_residue(int(a_coff), 2, p)

        div = N // r
        if div == 2 and quad_residue: continue
        if div == 4 and not quad_residue: continue
        return ((x, y), a_coff), points 
    
    return ((None, None), None), points


def draw_curve(x_coordinats, y_coordinats):
    _, ax = plt.subplots()
    ax.scatter(x_coordinats, y_coordinats)
    plt.show()

def gen_curve(l, m):
  while True:
    p = generate_prime(l)
    if p == None: return

    c_d_pair = div_to_squar(p, 1)
    if (c_d_pair == None): continue

    a, b = c_d_pair

    params = prov_t_cond(p, a, b, m)

    if len(params) == 0:
        print()
        continue
    
    for p, n_big, r in params:
        points = set()
        while len(points) != n_big:
            pair_p_a, points = generate_points(p, n_big, r, points)
            if pair_p_a != ((None, None), None):
                p_point, a_coff = pair_p_a
                if not check_infinite_point(n_big, *p_point, a_coff, p):
                    continue
                q_big = scalarMultiply(int(n_big // r), *p_point, p, a_coff)
                (p, a_coff, q_big, r)
                points_groupp = []

                for k in range(1, int(r) + 1):
                    new_point = scalarMultiply(k, q_big[0], q_big[1], p, a_coff)
                    points_groupp.append([new_point[0],new_point[1],  k])
                    
                random.shuffle(points_groupp )
                random_point = random.choice(points_groupp)
                return (p, a_coff, q_big, r, random_point)
            else:
                continue
    return

