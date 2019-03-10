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


# Returns byte in hex_string on byte_pos position. MSB has byte_pos = 0, LSB has byte_pos = -1.
def get_byte(hex_string, byte_pos):
    return [hex_string[i:i + 2] for i in range(0, len(hex_string), 2)][byte_pos]


def plot_byte_histogram(data, column_name, byte_pos):
    data[column_name] \
        .apply(lambda x: int(get_byte(x, byte_pos), 16)) \
        .plot.hist(bins=256)


def plot_msb_histogram(data, column_name):
    plot_byte_histogram(data, column_name, 0)


def plot_lsb_histogram(data, column_name):
    plot_byte_histogram(data, column_name, -1)


def plot_byte_heatmap(data, column_name, byte_pos):
    index = data['t1'].apply(lambda x: x // 1000000)
    columns = data[column_name].apply(lambda x: int(get_byte(x, byte_pos), 16))
    cross_table = pd.crosstab(
        index=index,
        columns=columns,
        dropna=False
    ).transpose().reindex(np.arange(0, 256)).transpose().fillna(0)

    plt.pcolor(cross_table)


def plot_msb_heatmap(data, column_name):
    plot_byte_heatmap(data, column_name, 0)


def plot_lsb_heatmap(data, column_name):
    plot_byte_heatmap(data, column_name, -1)


def generate_and_save_graphs(data, file_name_prefix, limit):
    data_size = len(data.index)

    plot_histogram_of_keygen_times(sampleData)
    plt.savefig(file_name_prefix + "_hist_keygen_times.png")
    plt.clf()

    if data_size <= limit:
        plot_keygen_time_for_generated_keys(sampleData)
        plt.savefig(file_name_prefix + "_keygen_times.png")
        plt.clf()

    plot_msb_histogram(sampleData, 'p')
    plt.savefig(file_name_prefix + "_hist_msb_p.png")
    plt.clf()

    plot_lsb_histogram(sampleData, 'p')
    plt.savefig(file_name_prefix + "_hist_lsb_p.png")
    plt.clf()

    plot_msb_heatmap(sampleData, 'p')
    plt.savefig(file_name_prefix + "_heatmap_msb_p.png")
    plt.clf()

    plot_lsb_heatmap(sampleData, 'p')
    plt.savefig(file_name_prefix + "_heatmap_lsb_p.png")
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
