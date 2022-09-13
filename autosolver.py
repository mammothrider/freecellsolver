# import pywin32
# import mouse
import contextlib
from tarfile import TarError
from tracemalloc import Statistic
import cv2
import os

class Card:
    def __init__(self) -> None:
        self.card = {}
        self.load_card()

    def load_card(self):
        path = os.path.dirname(__file__)
        image_folder = os.path.join(path, "image")
        for filename in os.listdir(image_folder):
            if not filename.endswith("jpg"):
                continue
            file = cv2.imread(os.path.join(image_folder, filename))
            number = filename.split(".")[0]
            self.card[int(number)] = file

    def find_card(self, target: cv2.Mat):
        result = 0
        tmpMin = 9999999
        for k, v in self.card.items():
            if v.rows == target.rows and v.cols == target.cols:
                # Calculate the L2 relative error between images.
                errorL2 = cv2.norm(v, target, cv2.CV_L2)
                # Convert to a reasonable scale, since L2 error is summed across all pixels of the image.
                similarity = errorL2 / (v.rows * v.cols)
                if similarity < tmpMin:
                    result = k
        return result


class GameStruct:
    def __init__(self) -> None:
        self.free = [0] * 4
        self.home = [0] * 4
        self.card = [[] for i in range(8)]


class Solver:

    def __init__(self) -> None:
        self.card = Card()
    
    def capture_window(self):
        pass
    
    def analyze(self, pic: cv2.Mat) -> GameStruct:
        game = GameStruct()
        rsz_img = cv2.resize(img, (1000, 650))
        gray = cv2.cvtColor(rsz_img, cv2.COLOR_BGR2GRAY)

        pixel = 30
        start = 200
        cut_height = 24
        cut_width = 16
        y = [82, 190, 299, 408, 516, 625, 734, 842]

        for i in range(7):
            for j in range(8):
                h = start + pixel * i
                w = y[j]
                if i*8+j>=52:
                    break
                card = gray[h:h+cut_height, w:w+cut_width]

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



    
    def call_solver(self, game: GameStruct):
        pass
    
    def solve(self, actions):
        pass

if __name__ == "__main__":
    # test
    img = cv2.imread("img.png")
    game = Solver().analyze(img)
    print(game)