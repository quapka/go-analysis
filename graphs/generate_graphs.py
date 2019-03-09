import pandas as pd
import matplotlib.pyplot as plt
import numpy as np


def show_histogram_of_keygen_times(data):
    data['t1'].plot.hist()
    plt.show()


def show_keygen_time_for_generated_keys(data):
    data.plot.bar(x='id', y='t1', width=1)

    num_of_rows = len(data.index)
    num_of_steps = num_of_rows // 10
    plt.xticks(np.arange(0, num_of_rows, num_of_steps), data.index.values[0::num_of_steps])

    plt.show()


def show_histogram_of_msb_of_prime_p(data):
    data['p']\
        .apply(lambda x: x[:2])\
        .apply(lambda x: int(x, 16))\
        .plot.hist(bins=256)

    plt.show()


def show_histogram_of_lsb_of_prime_p(data):
    data['p']\
        .apply(lambda x: x[-2:])\
        .apply(lambda x: int(x, 16))\
        .plot.hist(bins=256)

    plt.show()


sampleData = pd.read_csv('../sample-data/rsa_512b_100.csv', delimiter=';')

show_histogram_of_keygen_times(sampleData)
show_keygen_time_for_generated_keys(sampleData)
show_histogram_of_msb_of_prime_p(sampleData)
show_histogram_of_lsb_of_prime_p(sampleData)
