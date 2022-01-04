import fire


def main(
        #infile: str='test_input.txt',
        infile: str = 'input.txt',
):
    print('hi!')

    with open(infile, 'r') as f:
        lines = [
            int(x.strip()) for x in f.readlines()
        ]

    print(lines)

    p = lines[0]
    inc = 0
    for i in lines[1:]:
        if i > p:
            inc += 1
        p = i

    print(f'1: {inc}')

    inc = 0
    p = lines[0:3]
    s = sum(p)
    for i, n in enumerate(lines[3:]):
        p[i%3] = n
        s2 = sum(p)
        if s2 > s:
            inc += 1
        s = s2

    print(f'2: {inc}')


if __name__ == '__main__':
    fire.Fire(main)
