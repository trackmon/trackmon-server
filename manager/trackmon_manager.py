import sys

def main():
    if "-install" in sys.argv:
        print("Installing everything")
    elif "-installapi" in sys.argv:
        print("Installing API backend only")
    elif "-installdb" in sys.argv:
        print("Installing database only")
    elif "-installfrontend" in sys.argv:
        print("Installing frontend only")

if __name__ == "__main__":
    main()
