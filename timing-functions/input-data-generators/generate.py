import csv
import random


def generate_sbm_random(filename, num_of_values):
    """scalar_base_mult random"""

    with open(filename, 'w', newline='') as csvfile:
        for i in range(0, num_of_values):
            value = random.randrange(1, 2**256-432420386565659656852420866394968145599)
            value_string = format(value, 'X')
            label = "col_rnd_" + str(i)

            writer = csv.writer(csvfile, delimiter=';')
            writer.writerow([label, value_string])


generate_sbm_random("test.csv", 10)
