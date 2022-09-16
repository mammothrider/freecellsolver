# import pywin32
# import mouse
import time
import cv2
import os
import subprocess
import json

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

    def to_dict(self) -> dict:
        return dict(
            card = self.card,
            home = self.home,
            free = self.free,
        )


class Solver:

    def __init__(self) -> None:
        self.card = Card()
    
    def capture_window(self):
        pass
    
    def analyze(self, pic: cv2.Mat) -> GameStruct:
        game = GameStruct()
        rsz_img = cv2.resize(pic, (1000, 650))
        gray = cv2.cvtColor(rsz_img, cv2.COLOR_BGR2GRAY)

        # card 
        gap_height = 30
        start = 200
        cut_height = 24
        cut_width = 16
        y = [82, 190, 299, 408, 516, 625, 734, 842]

        for col in range(8):
            row = 0
            while True:
                h = start + gap_height * row
                w = y[col]
                card = gray[h:h+cut_height, w:w+cut_width]
                num = self.card.find_card(card)
                if num == 0:
                    break
                game.card[col].append(num)
                row += 1

        # cv2.rectangle(gray, (10, 20), (35, 60), color=(255, 0, 0))
        # cv2.rectangle(gray, (970, 625), (995, 645), color=(255, 0, 0))

        # free area
        free_start = 50
        free_height = 70
        space_width = 70
        gap = 40
        for i in range(4):
            free_left = free_start + (space_width + gap) * i
            card = gray[free_height:free_height+cut_height, free_left:free_left+cut_width]
            cv2.rectangle(gray, (free_left, free_height), (free_left + cut_width, free_height + cut_height), color=(255, 0, 0))
            num = self.card.find_card(card)
            game.free[i] = num

        home_start = 550
        for i in range(4):
            home_left = home_start + (space_width + gap) * i
            card = gray[free_height:free_height+cut_height, home_left:home_left+cut_width]
            cv2.rectangle(gray, (home_left, free_height), (home_left + cut_width, free_height + cut_height), color=(255, 0, 0))
            num = self.card.find_card(card)
            # home区颜色顺序和设计不对应，做调整
            game.home[num//100 - 1] = num

        return game

    def call_solver(self, game: GameStruct):
        data = json.dumps(game.to_dict())
        # print(data)
        result = subprocess.run(["./freecellsolver", data], capture_output=True, text=True)
        if result.stderr:
            print(result.stderr)
            return ""
        text = result.stdout
        # print(text)
        return json.loads(text)
    
    def solve(self, actions):
        pass

if __name__ == "__main__":
    # test
    # img = cv2.imread("img.png")
    # game = Solver().analyze(img)
    # print(game)
    game = GameStruct()
    game.card = [
        [411, 310, 108, 311, 203, 407, 403],
        [313, 106, 105, 408, 104, 410, 201],
        [306, 202, 204, 113, 401, 205, 307],
        [101, 405, 413, 102, 312, 309, 303],
        [302, 209, 208, 213, 409, 111],
        [304, 404, 206, 109, 412, 406],
        [301, 112, 402, 212, 210, 305],
        [211, 308, 107, 110, 207, 103],
    ]
    a = time.time()
    res = Solver().call_solver(game)
    print(time.time() - a)