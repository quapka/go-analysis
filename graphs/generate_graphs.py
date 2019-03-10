import pandas as pd
import matplotlib.pyplot as plt
import numpy as np


def plot_histogram_of_keygen_times(data):
    data['t1'].plot.hist()


# plotting of large amount of data takes a lot of time and memory
def plot_keygen_time_for_generated_keys(data):
    data.plot.bar(x='id', y='t1', width=1)

    num_of_rows = len(data.index)
    num_of_steps = num_of_rows // 10
    plt.xticks(np.arange(0, num_of_rows, num_of_steps), data.index.values[0::num_of_steps])


def plot_histogram_of_msb_of_prime_p(data):
    data['p']\
        .apply(lambda x: x[:2])\
        .apply(lambda x: int(x, 16))\
        .plot.hist(bins=256)


def plot_histogram_of_lsb_of_prime_p(data):
    data['p']\
        .apply(lambda x: x[-2:])\
        .apply(lambda x: int(x, 16))\
        .plot.hist(bins=256)


def generate_and_save_graphs(data, file_name_prefix, small_data_only):
    data_size = len(data.index)

    plot_histogram_of_keygen_times(sampleData)
    plt.savefig(file_name_prefix + "_hist_keygen_times.png")
    plt.clf()

    if data_size <= 10000 or not small_data_only:
        plot_keygen_time_for_generated_keys(sampleData)
        plt.savefig(file_name_prefix + "_keygen_times.png")
        plt.clf()

    plot_histogram_of_msb_of_prime_p(sampleData)
    plt.savefig(file_name_prefix + "_hist_msb_p.png")
    plt.clf()

    plot_histogram_of_lsb_of_prime_p(sampleData)
    plt.savefig(file_name_prefix + "_hist_lsb_p.png")
    plt.clf()


sampleData = pd.read_csv(
    '../sample-data/rsa_512b_10000.csv',
    delimiter=';',
    low_memory=False,
)


generate_and_save_graphs(sampleData, "rsa_512b_10000", True)