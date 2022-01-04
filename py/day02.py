import fire


def main(
        #infile: str='test_input.txt',
        infile: str = 'input.txt',
):
    print('hi!')

    with open(infile, 'r') as f:
        input = []
        for l in f.readlines():
            i = l.split(' ')
            i[0] = i[0][0]
            i[1] = int(i[1])
            input.append(tuple(i))

    pos = [0, 0]
    for d, x in input:
        if d == 'f':
            pos[0] += x
        elif d == 'd':
            pos[1] += x
        elif d == 'u':
            pos[1] -= x


    # print(input)

    print(f'2: {pos[0]*pos[1]}')

    pos = [0, 0]
    aim = 0
    for d, x in input:
        if d == 'f':
            pos[0] += x
            pos[1] += aim*x
        elif d == 'd':
            aim += x
        elif d == 'u':
            aim -= x


    # print(input)

    print(f'1: {pos[0]*pos[1]}')


if __name__ == '__main__':
    fire.Fire(main)
