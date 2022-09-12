import cv2

import numpy as np

name = [
    308, 209, 204, 101, 211, 107, 406, 313,
    109, 301, 413, 410, 312, 208, 412, 202,
    111, 113, 401, 407, 112, 405, 103, 210,
    403, 203, 212, 108, 310, 302, 201, 106,
    110, 105, 102, 404, 409, 206, 408, 304,
    205, 306, 213, 207, 309, 303, 311, 402,
    307, 411, 305, 104,
]

# load image
img = cv2.imread('img.png')

rsz_img = cv2.resize(img, (1000, 650)) # resize since image is huge

gray = cv2.cvtColor(rsz_img, cv2.COLOR_BGR2GRAY) # convert to grayscale

height, width, channels = img.shape

pixel = 30
start = 200
cut_height = 24
cut_width = 16
# start_width = 81
# gap = 92.8

y = [82, 190, 299, 408, 516, 625, 734, 842]

for i in range(7):
    for j in range(8):
        h = start + pixel * i
        w = y[j]
        if i*8+j>=52:
            break
        # crop = cv2.imwrite(f"image/{name[i*8+j]}.jpg", gray[h:h+cut_height, w:w+cut_width])
        # cv2.rectangle(gray, (w, h), (w + cut_width, h + cut_height), color=(0, 0, 0))

cv2.rectangle(gray, (10, 20), (35, 60), color=(255, 0, 0))
cv2.rectangle(gray, (970, 625), (995, 645), color=(255, 0, 0))

free_start = 50
home_start = 550
free_height = 70
space_width = 70
space_height = 100
gap = 40
for i in range(4):
    free_left = free_start + (space_width + gap) * i
    home_left = home_start + (space_width + gap) * i
    cv2.rectangle(gray, (free_left, free_height), (free_left + space_width, free_height + space_height), color=(255, 0, 0))
    cv2.rectangle(gray, (home_left, free_height), (home_left + space_width, free_height + space_height), color=(255, 0, 0))

cv2.imshow('img', gray)
cv2.waitKey(0)