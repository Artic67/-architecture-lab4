Лабораторна робота №4
---

### Виконавці:
* [Негруб Артем](https://github.com/Artic67)
* [Тарасюк Владислав](https://github.com/vtarasiuk)

### Package Loop:

Structure cmdQueue methods:
* push
* pull
* isEmpty

Structure EventLoop methods:
* init
* dispose
* Post
* AwaitClose
* Start
* isStopRequested
* isStopped
* stop
* listen
* verifyRunning

### Package Parser:

Contains:
* Structures PrintCmd and DeleteCmd with method Execute.
* cmdList - map of existing commands
* findCommand function to find command)
* Main function Parse to parse command and return it
* functions funcPrint and funcDelete to make sturncture PrintCmd and DeleteCmd from command

