from multiprocessing.sharedctypes import Value
import win32gui
import mouse
import time
import cv2
import os
import subprocess
import json
import numpy
from PIL import ImageGrab

import ctypes
 
ctypes.windll.shcore.SetProcessDpiAwareness(1)

class Card:
    def __init__(self) -> None:
        self.card = {}
        self.load_card()

    def load_card(self):
        path = os.path.dirname(__file__)
        image_folder = os.path.join(path, "image")
        for filename in os.listdir(image_folder):
            if not filename.endswith("png"):
                continue
            # file = cv2.imread(os.path.join(image_folder, filename), cv2.IMREAD_GRAYSCALE)
            file = cv2.imread(os.path.join(image_folder, filename), cv2.IMREAD_UNCHANGED)
            # file = cv2.cvtColor(file, cv2.COLOR_BGR2GRAY)
            # file = cv2.threshold(file, 120, 255, cv2.THRESH_BINARY)[1]

            number = filename.split(".")[0]
            self.card[int(number)] = file

    def find_card(self, target):
        result = 0
        tmpMax = 0
        # target = cv2.cvtColor(target, cv2.COLOR_BGR2GRAY)
        # target = cv2.threshold(target, 120, 255, cv2.THRESH_BINARY)[1]
        cache = {}
        for k, v in self.card.items():
            tmp = cv2.matchTemplate(target, v, cv2.TM_CCOEFF_NORMED)
            min_val, max_val, min_loc, max_loc = cv2.minMaxLoc(tmp)
            cache[k] = round(max_val, 3)
            if max_val > 0.7 and max_val > tmpMax:
                tmpMax = max_val
                result = k
        # print(result, cache)
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

    def print_game(self):
        for i in self.free + self.home:
            print(f"{i:4}", end="" )
        print()
        print("---------------------------------")
        count = 0
        while True:
            flag = True
            for card in self.card:
                if len(card) > count:
                    print(f"{card[count]:4}", end="")
                    flag = False
                else:
                    print(" ---", end="")

            print()
            count += 1
            if flag:
                break

    def check_legal(self):
        card_list = []
        for c in self.card:
            card_list.extend(c)
        card = set(card_list)
        target = set([(i + 1) * 100 + j + 1 for i in range(4) for j in range(13)])
        if len(card) == 52 and len(target - card) == 0:
            return

        not_found = list(target - card)
        raise ValueError(f"Unlegal {not_found}")


class Solver:

    def __init__(self) -> None:
        self.card = Card()
        self.bbox = (0, 0, 2000, 1300)

    def reset_window(self):
        title = "Solitaire Collection"
        hwnd = win32gui.FindWindow(None, title)
        if hwnd == 0:
            return
        win32gui.SetForegroundWindow(hwnd)
        # bbox = win32gui.GetWindowRect(hwnd)
        win32gui.MoveWindow(hwnd, self.bbox[0], self.bbox[1], self.bbox[2], self.bbox[3], True)
        return
    
    def capture_window(self):
        self.reset_window()
        img = ImageGrab.grab(self.bbox)
        img = numpy.array(img.convert("RGB"))
        img = cv2.cvtColor(img, cv2.COLOR_BGR2RGB)
        return img
    
    def analyze(self, pic, debug = False) -> GameStruct:
        game = GameStruct()
        # rsz_img = cv2.resize(pic, (1000, 650))
        # gray = cv2.cvtColor(rsz_img, cv2.COLOR_BGR2GRAY)
        target_img = pic

        # card 
        gap_height = 60
        start = 428
        cut_height = 54
        cut_width = 30
        y = [172, 388, 602, 817, 1032, 1246, 1462, 1676]

        for col in range(8):
            for row in range(8):
                if col + row * 8 >= 52:
                    continue
                h = int(start + gap_height * row)
                w = y[col]
                card = target_img[h:h+cut_height, w:w+cut_width]
                num = self.card.find_card(card)
                cv2.rectangle(target_img, (w, h), (w + cut_width, h + cut_height), color=(0, 0, 0))
                if num == 0:
                    break
                game.card[col].append(num)

        # free area
        # free_height = 90
        # free_col = [55, 163, 270, 378]
        # for i in range(4):
        #     free_left = free_col[i]
        #     card = target_img[free_height:free_height+cut_height, free_left:free_left+cut_width]
        #     cv2.rectangle(target_img, (free_left, free_height), (free_left + cut_width, free_height + cut_height), color=(0, 0, 0))
        #     num = self.card.find_card(card)
        #     game.free[i] = num

        # home_col = [546, 654, 763, 871]
        # for i in range(4):
        #     home_left = home_col[i]
        #     card = target_img[free_height:free_height+cut_height, home_left:home_left+cut_width]
        #     cv2.rectangle(target_img, (home_left, free_height), (home_left + cut_width, free_height + cut_height), color=(0, 0, 0))
        #     num = self.card.find_card(card)
        #     # home区颜色顺序和设计不对应，做调整
        #     color = (i + 1) % 4
        #     game.home[color] = num

        if debug:
            game.print_game()
            cv2.imshow("", target_img)
            cv2.waitKey(0)

        return game

    def call_solver(self, game: GameStruct):
        data = json.dumps(game.to_dict())
        # print(data)
        result = subprocess.run(["./freecellsolver", data], capture_output=True, text=True)
        if result.stderr:
            print(result.stderr)
            return ""
        text = result.stdout
        # print(result.stdout, result.stderr)
        return json.loads(text)
    
    def get_pos(self, action):
        gap_height = 29.8
        row_start = 213
        card_width = 75
        col_y = [86, 193, 301, 408, 516, 623, 731, 838]

        free_height = 90
        free_col = [55, 163, 270, 378]

        start_pos = end_pos = [0, 0]
        act = action["Action"]
        if act == "Up" or act == "Move" or act == "Free":
            x = col_y[action["FCol"]] + card_width/2
            y = row_start + gap_height * min(action["FRow"], 10) + gap_height
            start_pos = [int(x * 2), int(y * 2)]
        elif act == "FreeMove":
            x = free_col[action["FCol"]] + card_width/2
            y = free_height + gap_height
            start_pos = [int(x * 2), int(y * 2)]

        if act == "Move" or act == "FreeMove":
            x = col_y[action["TCol"]] + card_width/2
            y = row_start + gap_height * min(action["TRow"], 10) + gap_height
            end_pos = [int(x * 2), int(y * 2)]
        elif act == "Free":
            x = free_col[action["TCol"]] + card_width/2
            y = free_height + gap_height
            end_pos = [int(x * 2), int(y * 2)]

        return start_pos, end_pos

    def solve(self, actions):
        if not actions:
            return
        self.reset_window()
        stop = False
        def callback():
            nonlocal stop
            stop = True
            print("STOP")

        mouse.on_button(callback, buttons=(mouse.RIGHT), types=(mouse.UP))
        for index, action in enumerate(actions):
            start_pos, end_pos = self.get_pos(action)
            act = action["Action"]
            if act == "Up" and index == 0:
                mouse.move(start_pos[0], start_pos[1], True, 0.2)
                mouse.double_click()
            elif act == "Free" or act == "Move" or act == "FreeMove":
                # mouse.drag(start_pos[0], start_pos[1], end_pos[0], end_pos[1], True, 0.5)
                mouse.move(start_pos[0], start_pos[1], True, 0.2)
                mouse.click()
                mouse.move(end_pos[0], end_pos[1], True, 0.2)
                mouse.click()
            elif act == "Up":
                time.sleep(0.5)

            if stop:
                break
            print(index, action, start_pos, end_pos)
            time.sleep(0.5)
    

if __name__ == "__main__":
    solver = Solver()
    image = Solver().capture_window()
    game = Solver().analyze(image)
    game.print_game()
    game.check_legal()
    print("Call solver")
    actions = solver.call_solver(game)
    if not actions:
        data = game.to_dict()
        print("Failed", data)
    solver.solve(actions)

    # test capture
    # image = Solver().capture_window()
    # game = Solver().analyze(image, debug = True)
    # game.check_legal()

    # test solver
    # game = GameStruct()
    # a = [
    #     103, 403, 109, 204, 211, 112, 404, 207,
    #     405, 402, 205, 108, 206, 309, 212, 406,
    #     104, 110, 306, 201, 308, 407, 102, 303,
    #     304, 107, 310, 111, 203, 307, 313, 209,
    #     411, 412, 408, 208, 409, 302, 301, 312,
    #     413, 213, 311, 410, 401, 105, 305, 210,
    #     113, 202, 101, 106,
    # ]
    # for i in range(8):
    #     game.card[i] = [a[j] for j in range(i, 52, 8)]
    # game.print_game()
    # actions = solver.call_solver(game)
    # print(actions)


    # action = []
    # with open("test_actions.txt") as file:
    #     text = file.read()
    #     action = json.loads(text)
    # solver.solve(action)

    # img = solver.capture_window()
    # img = cv2.imread("img3.png")
    # game = solver.analyze(img)
    # actions = solver.call_solver(game)
    # solver.solve(actions)