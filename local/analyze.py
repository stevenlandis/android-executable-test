import json


def main():
    objs = []
    with open("local/log.txt", "r") as f:
        for line in f:
            objs.append(json.loads(line))
    diffs = []
    for i in range(len(objs) - 1):
        t0 = objs[i]["t"]
        t1 = objs[i + 1]["t"]
        diffs.append(t1 - t0)
    diffs.sort()
    print(diffs)


if __name__ == "__main__":
    main()
