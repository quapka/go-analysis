import csv
import random
import modular


def generate_random_point_by_weight(weight):
    a = -3
    b = 41058363725152142129326129780047268409114441015993725554835256314039467401291
    p = 2 ** 256 - 2 ** 224 + 2 ** 192 + 2 ** 96 - 1

    while True:
        x = generate_random_by_weight(256, weight)
        y_2 = x ** 3 + a * x + b
        y = modular.modular_sqrt(y_2, p)

        if y != 0:
            return x, y


def generate_random_by_weight_range(bit_size, min_weight, max_weight):
    return generate_random_by_weight(bit_size, random.randrange(min_weight, max_weight + 1))


def generate_random_by_weight(bit_size, weight):
    result = 0
    for power in random.sample(range(0, bit_size), weight):
        result += 2**power
    return result


def generate_random_by_weight_to_file(filename, num_of_values, bit_size, weight):
    with open(filename, 'w', newline='') as csvfile:
        for i in range(0, num_of_values):
            value = generate_random_by_weight(bit_size, weight)
            value_string = format(value, 'X')
            label = "col_rnd_weight_" + weight + "_" + str(i)

            writer = csv.writer(csvfile, delimiter=';')
            writer.writerow([label, value_string])


def generate_random_to_file(filename, num_of_values, min_value, max_value):
    with open(filename, 'w', newline='') as csvfile:
        for i in range(0, num_of_values):
            value = random.randrange(min_value, max_value)
            value_string = format(value, 'X')
            label = "col_rnd_" + str(i)

            writer = csv.writer(csvfile, delimiter=';')
            writer.writerow([label, value_string])
