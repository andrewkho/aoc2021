from dataclasses import dataclass
from typing import *
import numpy as np
import itertools

import fire


def main(
        #infile: str='test_input.txt',
        infile: str = 'input.txt',
):
    print('hi!')

    with open(infile, 'r') as f:
        lines = [
            x for x in f.readlines()
        ]

    nums = [int(x) for x in lines[0].split(',')]
    #print(nums)
    i = 2
    boards = []
    while i < len(lines):
        b = []
        for l in range(5):
            b.append(read_bingo_line(lines[i+l]))
        boards.append(Board(np.array(b)))
        i += 6

    for b in boards:
        print(b)

    winners = np.zeros(len(boards), dtype=np.int32)
    for n in nums:
        for i, b in enumerate(boards):
            if winners[i] == 1:
                continue
            if b.mark(n) and b.check():
                print(f'Winner: {(b.unmarked_sum(), n, b.unmarked_sum()*n)}')
                winners[i] = 1
                if winners.sum() == len(boards):
                    return


def read_bingo_line(line) -> List[int]:
    l = []
    for x in range(5):
        s = x*3
        e = x*3+2
        l.append(
            int(line[s:e])
        )

    return l


@dataclass
class Board:
    b: np.array
    N: int = 5

    def __post_init__(self):
        pass
        self.lu = dict()
        self.m = np.zeros(self.b.shape, dtype=np.int32)
        for i, j in itertools.product(range(self.N), range(self.N)):
            self.lu[self.b[i, j]] = (i, j)

    def mark(self, num: int) -> bool:
        if num in self.lu:
            idx = self.lu[num]
            self.m[idx] = 1
            return True
        return False

    def check(self) -> bool:
        for i in range(self.N):
            if self.m[i, :].sum() == 5:
                print(self, self.m)
                return True
            if self.m[:, i].sum() == 5:
                print(self, self.m)
                return True

        return False

    def unmarked_sum(self) -> int:
        return np.multiply(1-self.m, self.b).sum()


if __name__ == '__main__':
    fire.Fire(main)
