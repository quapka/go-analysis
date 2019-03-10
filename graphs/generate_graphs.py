import pandas as pd
import matplotlib.pyplot as plt
import numpy as np
import argparse
import os


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


def generate_and_save_graphs(data, file_name_prefix, limit):
    data_size = len(data.index)

    plot_histogram_of_keygen_times(sampleData)
    plt.savefig(file_name_prefix + "_hist_keygen_times.png")
    plt.clf()

    if data_size <= limit:
        plot_keygen_time_for_generated_keys(sampleData)
        plt.savefig(file_name_prefix + "_keygen_times.png")
        plt.clf()

    plot_histogram_of_msb_of_prime_p(sampleData)
    plt.savefig(file_name_prefix + "_hist_msb_p.png")
    plt.clf()

    plot_histogram_of_lsb_of_prime_p(sampleData)
    plt.savefig(file_name_prefix + "_hist_lsb_p.png")
    plt.clf()


# parse arguments
parser = argparse.ArgumentParser(description="Graphs generator.")
parser.add_argument('file', type=argparse.FileType('r'), nargs=1)
parser.add_argument('-n', default='', metavar='name_prefix',
                    help='Prefix of the filenames of the generated graphs. Default is the data filename.')
parser.add_argument('-l', metavar='limit', type=int, default=10000,
                    help='Some graphs can take only small amount of data. '
                         'This specifies the upper limit. The default is 10000.')
args = parser.parse_args()

# load data from csv
sampleData = pd.read_csv(
    args.file[0],
    delimiter=';',
    low_memory=False,
)

# generate graphs
file_basename = os.path.splitext(os.path.basename(args.file[0].name))[0]
prefix = file_basename if args.n == '' else args.n
generate_and_save_graphs(sampleData, prefix, args.l)
