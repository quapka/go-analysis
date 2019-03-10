import pandas as pd
import matplotlib.pyplot as plt
import argparse
import os
import plot_functions as pf


def rsa_generate_and_save_graphs(data, file_name_prefix, limit):
    data_size = len(data.index)

    pf.plot_histogram_of_keygen_times(sampleData)
    plt.savefig(file_name_prefix + "_hist_keygen_times.png")
    plt.clf()

    if data_size <= limit:
        pf.plot_keygen_time_for_generated_keys(sampleData)
        plt.savefig(file_name_prefix + "_keygen_times.png")
        plt.clf()

    pf.plot_msb_histogram(sampleData, 'p')
    plt.savefig(file_name_prefix + "_hist_msb_p.png")
    plt.clf()

    pf.plot_lsb_histogram(sampleData, 'p')
    plt.savefig(file_name_prefix + "_hist_lsb_p.png")
    plt.clf()

    pf.plot_msb_heatmap(sampleData, 'p')
    plt.savefig(file_name_prefix + "_heatmap_msb_p.png")
    plt.clf()

    pf.plot_lsb_heatmap(sampleData, 'p')
    plt.savefig(file_name_prefix + "_heatmap_lsb_p.png")
    plt.clf()


# parse arguments
parser = argparse.ArgumentParser(description="Graphs generator.")
parser.add_argument('file', type=argparse.FileType('r'), nargs=1)
parser.add_argument('-m', choices=['rsa', 'ecc'], required=True,
                    help='Modes of generating.')
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
rsa_generate_and_save_graphs(sampleData, prefix, args.l)
