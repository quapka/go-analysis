setwd('/home/qup/go/src/github.com/quapka/go-analysis/timing-functions/timers/src/')
#data <- read.csv('elliptic_ScalarBaseMult/data.csv', sep=';', header=TRUE)
#means <- apply(data, 2, mean)
#plot(means)

#meds <- apply(data, 2, median)
#plot(meds)

data <- read.csv('elliptic_Add/data.csv', sep=';', header=TRUE)
sbmData <- read.csv('elliptic_ScalarBaseMult/data_weights.csv', sep=';', header=TRUE)
dData <-read.csv('elliptic_Double/data.csv', sep = ';', header = T)

# means <- apply(data, 2, mean)
# plot(means, pch=16)
# plot(data$col_rnd_points_weight1_1weight2_2_0)
# outliers <- boxplot(data$col_rnd_points_weight1_1weight2_2_0)$out

# col <- data[-which(data$col_rnd_points_weight1_1weight2_2_0 %in% outliers),]

# meds <- apply(data, 2, median)
# plot(meds)

# stds <- apply(data, 2, sd)
# plot(stds)
library(gplots)
removeOutliers <- function(col) {
  outliers <- boxplot(plot=FALSE, col)$out
  col <- col[-which(col %in% outliers)]
  # av <- mean(col)
  return(col)
}

dataWithoutOutliers <- apply(data, 2, removeOutliers)
means <- lapply(dataWithoutOutliers, mean)
m <- matrix(unlist(means), nrow=256)

meds <- lapply(dataWithoutOutliers, median)
m2 <- matrix(unlist(meds), nrow=256)

sbmWO <- apply(sbmData, 2, removeOutliers)
