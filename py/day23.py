import time
import heapq
import math
from dataclasses import dataclass, field
from collections import deque, defaultdict, Counter
from typing import *

from scipy.signal import convolve2d

import numpy as np
import itertools
import re

import fire

def main(
        mode: str = 'auto',
        moves: Optional[str] = None,
):
    print('hi!')

    board = Board()
    board.print()

    if mode == 'interactive':
        interactive(board)
    elif mode == 'playback':
        playback(board, moves)
    else:
        t0 = time.time()
        best_score, best_path = solve_best(board)
        print('2:', best_score, time.time() - t0, 's', best_path)


def playback(board, moves):
    for start, end in moves:
        board.move_pos(start, end)
        board.print()
        time.sleep(0.25)


def interactive(board):
    history = []
    while True:
        try:
            valid_moves = [(s[1], e[1]) for s, e in board.get_valid_moves()]
            print(f'valid_moves:', valid_moves)
            i = input('from to: ')
            if i == 'z':
                f, t = history.pop()
                board.move(t, f, undo=-1)
            else:
                f, t = map(int, i.strip().split(' '))
                if (f, t) not in valid_moves:
                    print("Invalid move!", (f, t))
                    continue
                board.move(f, t)
                history.append((f, t))
            print(f'score: {board.score}')
            board.print()
        except Exception as e:
            print("Exception!")
            print(str(e))
            board.print()


def solve_best(board):
    best_score = np.inf
    best_path = None
    history = []
    seen = dict()
    cache_hits = 0

    def recurse():
        nonlocal best_score, best_path, cache_hits

        state = board.get_str()
        if state in seen and seen[state] <= board.score:
            cache_hits += 1
            return
        seen[state] = board.score

        if board.winner():
            if board.score < best_score:
                best_score = min(best_score, board.score)
                best_path = history.copy()
                print(best_score, board.score, best_path)
            return

        for start, end in board.get_valid_moves():
            board.move_pos(start, end)
            history.append((start, end))
            recurse()
            f, t = history.pop()
            board.move_pos(t, f, undo=-1)

    recurse()

    print(f'cache_hits:', cache_hits, 'cache_size:', len(seen))
    return best_score, best_path


class Board:
    CORRECT_COLS = {
        'A': 2,
        'B': 4,
        'C': 6,
        'D': 8
    }
    ENERGY = {
        'A': 1,
        'B': 10,
        'C': 100,
        'D': 1000,
    }

    def __init__(self):
        self.score = 0
        self.b = ['.' for _ in range(11)]
        self.b[2] += 'DDDB'
        self.b[4] += 'DCBA'
        self.b[6] += 'CBAA'
        self.b[8] += 'BACC'

    def get_str(self):
        return ''.join(col for col in self.b)

    def winner(self):
        return (self.b[2][1:] == 'AAAA' and
                self.b[4][1:] == 'BBBB' and
                self.b[6][1:] == 'CCCC' and
                self.b[8][1:] == 'DDDD')

    def print(self):
        lines = []
        lines.append('  01234567890')
        lines.append(' ' + '#'*13)
        # row 1
        lines.append('0#' + ''.join(self.b[c][0] for c in range(11)) + '#')
        lines.append('1###' + '#'.join(self.b[c][1] for c in (2, 4, 6, 8)) + '###')
        for row in range(2, 5):
            lines.append(f'{row}  #' + '#'.join(self.b[c][row] for c in (2, 4, 6, 8)) + '#  ')
        lines.append('   ' + '#'*9)
        print('\n'.join(lines))

    def col_to_pos(self, col, strt: int):
        row = 0
        if col in [2, 4, 6, 8]:
            for row, v in enumerate(self.b[col]):
                if v in 'ABCD':
                    row -= (1-strt)
                    break
        return row, col

    def move(self, start_col, end_col, undo=1):
        start = self.col_to_pos(start_col, 1)
        end = self.col_to_pos(end_col, 0)

        self.move_pos(start, end, undo)

    def move_pos(self, start, end, undo=1):
        c = self.b[start[1]][start[0]]
        self.b[end[1]] = self.b[end[1]][:end[0]] + c + self.b[end[1]][end[0]+1:]
        self.b[start[1]] = self.b[start[1]][:start[0]] + '.' + self.b[start[1]][start[0]+1:]

        self.score += undo * self.ENERGY[c]*(start[0] + end[0] + abs(start[1] - end[1]))

    def get_valid_moves(self):
        valid = []
        for start_col in range(11):
            start = self.col_to_pos(start_col, 1)
            c = self.b[start[1]][start[0]]
            if c not in 'ABCD':
                continue

            # If start is in correct hole, then don't move
            if start[0] > 0 and start[1] == self.CORRECT_COLS[c]:
                if self.b[start[1]][start[0]+1:] == c*(5 - (start[0]+1)):
                    continue

            start_valids = list()
            # go left to right
            for end_col in range(start_col+1, 11):
                end = self.col_to_pos(end_col, 0)
                if self.b[end[1]][end[0]] != '.':
                    # path is blocked from here on out
                    break

                if self.move_is_final(start, end):
                    return [(start, end)]
                elif self.move_is_valid(start, end):
                    start_valids.append((start, end))

            # go right to left
            for end_col in range(start_col-1, -1, -1):
                end = self.col_to_pos(end_col, 0)
                if self.b[end[1]][end[0]] != '.':
                    # path is blocked from here on out
                    break

                if self.move_is_final(start, end):
                    return [(start, end)]
                elif self.move_is_valid(start, end):
                    start_valids.append((start, end))

            valid.extend(start_valids)

        return valid

    def move_is_final(self, start, end) -> bool:
        if end[0] == 0:
            return False

        # the end must be in the correct col
        c = self.b[start[1]][start[0]]
        correct_col = self.CORRECT_COLS[c]
        if end[1] != correct_col:
            return False

        # Check every value below to see if it matches expected char
        return self.b[end[1]][end[0]+1:] == c*(5 - (end[0]+1))

    def move_is_valid(self, start, end) -> bool:
        # we don't check this unless we're definitely NOT ending in the correct hole
        if start[0] == 0 and end[0] == 0:
            return False

        if start[0] > 0 and end[0] != 0:
            return end[0] == 0

        if end[1] in (2, 4, 6, 8):
            return False

        c = self.b[start[1]][start[0]]
        if c not in 'ABCD':
            return False

        return True


if __name__ == '__main__':
    fire.Fire(main)
