import cv2

import numpy as np

name = [
    103, 403, 109, 204, 211, 112, 404, 207,
    405, 402, 205, 108, 206, 309, 212, 406,
    104, 110, 306, 201, 308, 407, 102, 303,
    304, 107, 310, 111, 203, 307, 313, 209,
    411, 412, 408, 208, 409, 302, 301, 312,
    413, 213, 311, 410, 401, 105, 305, 210,
    113, 202, 101, 106,
]

# load image
img = cv2.imread('img3.png')
print(img.shape)

# rsz_img = cv2.resize(img, (1000, 650)) # resize since image is huge

# gray = cv2.cvtColor(rsz_img, cv2.COLOR_BGR2GRAY) # convert to grayscale
target_img = img

height, width, channels = img.shape

pixel = 60
start = 428
cut_height = 54
cut_width = 30
# start_width = 81
# gap = 92.8

y = [172, 388, 602, 816, 1032, 1246, 1462, 1676]

def clean_img(pic):
    tmp_gray = cv2.cvtColor(pic, cv2.COLOR_BGR2GRAY)
    blur = cv2.GaussianBlur(tmp_gray, (3,3), 0)
    thresh = cv2.threshold(blur, 200, 255, cv2.THRESH_BINARY_INV + cv2.THRESH_OTSU)[1]
    rect = cv2.boundingRect(thresh)
    crop = pic[rect[1]:rect[1] + rect[3], rect[0]:rect[0] + rect[2]]
    return crop, rect

for i in range(7):
    for j in range(8):
        h = start + pixel * i
        h = int(h)
        w = y[j]
        if i*8+j>=52:
            break

        pic = target_img[h:h+cut_height, w:w+cut_width]
        crop, rect = clean_img(pic)

        crop = cv2.imwrite(f"image/{name[i*8+j]}.png", crop)
        cv2.rectangle(target_img, (w + rect[0], h + rect[1]), (w + rect[0] + rect[2], h + rect[1] + rect[3]), color=(0, 0, 0))

# cv2.rectangle(gray, (10, 20), (35, 60), color=(255, 0, 0))
# cv2.rectangle(gray, (970, 625), (995, 645), color=(255, 0, 0))

# free_start = 50
# home_start = 550
# free_height = 70
# space_width = 70
# space_height = 100
# gap = 40
# for i in range(4):
#     free_left = free_start + (space_width + gap) * i
#     home_left = home_start + (space_width + gap) * i
#     cv2.rectangle(gray, (free_left, free_height), (free_left + space_width, free_height + space_height), color=(255, 0, 0))
#     cv2.rectangle(gray, (home_left, free_height), (home_left + space_width, free_height + space_height), color=(255, 0, 0))

cv2.imshow('img', target_img)
cv2.waitKey(0)