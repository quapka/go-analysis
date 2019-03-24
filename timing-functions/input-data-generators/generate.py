import csv
import random


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
            label = "col_rnd_" + str(i)

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


# generate_random_to_file("test.csv", 10, 1, 2**256-432420386565659656852420866394968145599)
generate_random_by_weight_to_file("test.csv", 10, 256, 200)
