import pandas as pd
import matplotlib.pyplot as plt
import numpy as np


def plot_histogram_of_keygen_times(data):
    data['t1'].plot.hist(bins=100)


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