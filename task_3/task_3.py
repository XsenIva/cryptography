import sympy
import random
import genEllptic

f = 'ell_curve.txt'
b = 'bank.txt' 
c = 'client.txt'

message = ""
alpha = 0

def function(x, y):
  return x + y 


def continue_or_break():
   ans = input("!! Продолжить выполнение? !!")
   _ = ans

def trunc_(file_name):
  with open(file_name, 'w') as f:
    f.truncate(0)

def close_(file_name):
  with open(file_name, 'w') as f:
    f.close()

def point_on_curve(x, y, p, a):
  flag = True
  if (y**2) % p != (x**3 + a * x ) % p:
    flag = False
    print(f"Точка {x}, {y} не лежит кривой, безопасность нарушена!!!")
  return flag

def clean_string(s):
    str = ''.join(char for char in s if char.isdigit())
    return int(str)

def clean_coords(s):
    cleaned = ''.join(char for char in s if char.isdigit() or char.isspace())
    return [num for num in cleaned.split() if num]


def write_curve( p, a_coff, q_big, r, random_point, message):
  with open(f, 'a+') as fi:
    fi.write("\nОткрытый ключ банка\n")
    fi.write("P " + str(p) + "\n")
    fi.write("A " + str(a_coff) + "\n")
    fi.write("q " + str( q_big[0]) + " " + str( q_big[1]) + "\n")
    fi.write("r " + str(r) + "\n")
    fi.write("t " + str(random_point[0]) + " " + str(random_point[1]) + " " + str(random_point[2]) + "\n")
    fi.write("\nОткрытый ключ клиента\n")
    fi.write("P " + str(p) + "\n")
    fi.write("A " +  str(a_coff) + "\n")
    fi.write("q " + str( q_big[0]) + " " + str( q_big[1]) + "\n")
    fi.write("r " + str(r) + "\n")
    fi.write("t " + str(random_point[0]) + " " + str(random_point[1]) + "\n")
    fi.write("(" + message + ")")

  print("\n" + f"Сгенерирована кривая (P, A ,Q, r) = " + f"{p}, " + f"{a_coff}, " + f"{q_big}, " + f"{r}" + "\n" )
  print("Открытый ключ клиента: ")
  print("(P, A ,Q, r) = "+ f"{p}, " + f"{a_coff}, " + f"{q_big}, " + f"{r}")
  print("точка P из <Q> = "+ str(random_point[0]) + ", " + str(random_point[1]) + "\n")

  print("Открытый ключ банка:")
  print("(P, A ,Q, r) = "+ f"{p}, " + f"{a_coff}, " + f"{q_big}, " + f"{r}")
  print("точка P из <Q> = "+ str(random_point[0]) + ", " + str(random_point[1]))
  print("секретное l =  " + str(random_point[2]))


def step_first():
  print("\n--Генерация R'--")
  with open(f, 'r') as fi:
     lines = fi.readlines()
     for l in lines:
       l = l.split(" ")

  p, a_coff, r = clean_string(lines[2]), clean_string(lines[3]), clean_string(lines[5][:-2])
  if not sympy.isprime(int(r)):
    print("Порядок группы r не простой")

  coords_q = clean_coords(lines[4])
  q_big = [int(coords_q[0]), int(coords_q[1])]
  
  point_r = 0, 0
  while function(int(point_r[0]), int(point_r[1])) == 0:
    k = random.randint(1, r)
    print("сгенерировалось k' = ", k)
    point_r = genEllptic.scalarMultiply(k, q_big[0], q_big[1], p, a_coff)
   
    if point_r[0] == None or point_r[0] == None :
      continue
    if not point_on_curve(int(point_r[0]), int(point_r[1]), p, a_coff):
      return
    
    print("точка R' = ", point_r)

    with open(b, 'a+') as bi:
      bi.write("k' = " +  str(k) + "\n")
      bi.write("R' = "+ str(point_r[0]) +", "+ str(point_r[1]))
    print("\nБанк --R'--> Клиент \n")

def step_second():
  print("\n--Генерация R --")
  with open(f, 'r') as fi:
     lines = fi.readlines()
     for l in lines:
       l = l.split(" ")

  p, a_coff, r = clean_string(lines[9]), clean_string(lines[10]), clean_string(lines[12][:-2])
  mess = lines[14]
  mess = str(mess[1:-1])

  if not sympy.isprime(int(r)):
    print("Порядок группы r не простой")

  with open(b, 'r') as bi:
    line_r = bi.readlines()[-1]
    bank_r = line_r[5:].split(",")
 

  point_r = 0, 0
  while function(int(point_r[0]), int(point_r[1])) == 0:
    global alpha
    alpha = random.randint(1, r)
    print("сгенерировалось alpha = ", alpha)

    bank_x, bank_y = int(bank_r[0]), int(bank_r[1])

    point_r = genEllptic.scalarMultiply(alpha, bank_x, bank_y, p, a_coff)
    
    if point_r[0] == None or point_r[0] == None :
      continue
    print("точка R = ", point_r)
    
    if not point_on_curve(int(point_r[0]), int(point_r[1]), p, a_coff):
      return 

    betta = (function(int(point_r[0]), int(point_r[1])) * int(sympy.mod_inverse(int(function(bank_x, bank_y)), int(r)))) % r
    print("сгенерировалось betta = ", betta)
     

    m_new = sympy.mod_inverse(alpha, int(r))* betta *hash(mess) % r
    print("m' = ", m_new )

    with open(c, 'a+') as ci:
      ci.write("R = "+ str(point_r[0]) +", "+ str(point_r[1]))
      ci.write("\nm' = " + str(m_new))

    print("\nБанк <--m'-- Клиент \n")


def step_third():
  print("\n--Генерация s' --")
  with open(f, 'r') as fi:
     lines = fi.readlines()
     for li in lines:
       li = li.split(" ")
  
  r,  l_str = clean_string(lines[5][:-2]), lines[6].split(" ")
  l = int(l_str[-1])

  with open(b, 'r') as bi:
    r_lines = bi.readlines()
    line_k =  r_lines[-2][5:]
    line_r =  r_lines[-1][5:].split(",")
  
  with open(c, 'r') as ci:
    r_lines = ci.readlines()
    line_mess = r_lines[-1][5:].split(".")[0]
  
  rx, ry = int(line_r[0]), int(line_r[1])
  s_ = int(line_k) + l * function(rx ,ry) * int(line_mess) % r
  with open(b, 'a+') as bi: 
    bi.write("\ns' = " + str(s_))

  print("s' = ", str(s_))
  print("\nБанк --s'--> Клиент \n")
  

def step_fourth():
  global alpha
  print("\n--Проверка равенства --")
  
  with open(f, 'r') as fi:
     lines = fi.readlines()
     for l in lines:
       l = l.split(" ")

  p, a_coff, r = clean_string(lines[9]), clean_string(lines[10]), clean_string(lines[12][:-2])

  if not sympy.isprime(int(r)):
    print("Порядок группы r не простой")
  
  coords_q = clean_coords(lines[11])
  q_big = [int(coords_q[0]), int(coords_q[1])]

  with open(b, 'r') as bi:
    r_lines = bi.readlines()
    line_s = r_lines[-1][5:].split(".")[0]
    line_r = r_lines[-2]
    bank_r = line_r[5:].split(",")
    
  with open(f, 'r') as fi:
    r_lines = fi.readlines()
    coords_p = clean_coords(r_lines[13])

  with open(c, 'r') as ci:
    r_lines = ci.readlines()
    mess = r_lines[-1][5:]

  bank_x, bank_y = int(bank_r[0]), int(bank_r[1])
 
  left = genEllptic.scalarMultiply(int(line_s), q_big[0], q_big[1], p, a_coff)

  p_funk = genEllptic.scalarMultiply(function(bank_x, bank_y), int(coords_p[0]), int(coords_p[1]), p, a_coff)
  right_mult = genEllptic.scalarMultiply(int(mess), p_funk[0], p_funk[1], p, a_coff)
  right = genEllptic.add_ellip_points(bank_x, bank_y, right_mult[0], right_mult[1], a_coff, p)
  print(left,  "=",  right)
  if right != left:
    print("Подпись s' не действительна ")
    return None
  else:
    print("Подпись s' действительна ")
    s = alpha * int(line_s) % r
    make_coin(s)



def repayment_coin():
  print("\n--Проверка магазина --")
  with open(c, 'r') as ci:
    r_lines = ci.readlines()
    line_r = r_lines[-1][4:].split(", ")
  
  with open(f, 'r') as fi:
     lines = fi.readlines()
     for l in lines:
       l = l.split(" ")

  p, a_coff, r = clean_string(lines[9]), clean_string(lines[10]), clean_string(lines[12][:-2])
  px, py = clean_coords(lines[-1][13])
  
  if not sympy.isprime(int(r)):
    print("Порядок группы r не простой")
  
  coords_q = clean_coords(lines[11])
  line_q = [int(coords_q[0]), int(coords_q[1])]
  
  if line_r[0] == "" or  line_r[0] == " " or hash(line_r[0]) == 0:
     print("Подпись s' не действительна, равна 0")
     
  if function(int(line_r[1]), int(line_r[2])) == 0:
    print("Подпись s' не действительна f(подпись) = 0 ")
  
  m_fr = function(int(line_r[1]), int(line_r[2])) * hash(line_r[0]) % int(r)
  left = genEllptic.scalarMultiply(int(line_r[3]), line_q[0], line_q[1],  p, a_coff)

  p_funk = genEllptic.scalarMultiply(int(m_fr) , int(px), int(py), p, a_coff)
  rigth = genEllptic.add_ellip_points(int(line_r[1]),int(line_r[2]), p_funk[0],p_funk[1],  a_coff, p)
  
  if rigth != left:
    print("Подпись s' не действительна")
  else:
    print("Подпись s' действительна")


def to_step():
   global message
   param = int(input("На какой шаг вернуться?"))
   if param == 1:
     step_first()
     to_step()
   else: 
     if param == 2 :
      step_first()
      to_step()
     else:
       if param == 3 :
        step_third()
        to_step()
       else:
        if param == 4:
          step_fourth()
          to_step()
        else:
          repayment_coin()
          to_step()
   
def make_coin(ss):
  with open(c, 'r') as ci:
    r_lines = ci.readlines()
    line_r = clean_coords(r_lines[-2][4:])
  
  if ss != None:
    print("\nЭлектронная монета (m, R, s) = ", message, "(", line_r[0], line_r[1], ")", ss, "\n")
  
  with open(c, 'a+') as ci:
    ci.write(f"\nc = {message}, {str(int(line_r[0]))}, {str(int(line_r[1]))}, {str(ss)} ")



def main(): 
  trunc_(f)
  trunc_(c)
  trunc_(b)
  print("\n" + "Ввод параметров для генерации эллиптической кривой")
  l = int(input("l = "))
  m = int(input("m = "))
  if l < 4 : 
    print("Маленькое l")
    return
  global message
  message = input("\n" + "Введите сообщение m: ")

  p, a_coff, q_big, r, random_point = genEllptic.gen_curve(l, m)
  write_curve(p, a_coff, q_big, r, random_point, message)

  step_first()
  continue_or_break()

  step_second()
  continue_or_break()
  
  step_third()
  continue_or_break()

  step_fourth()
  continue_or_break()
  
  repayment_coin()
  to_step()

main()
