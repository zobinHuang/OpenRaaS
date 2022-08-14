// 该脚本基于 docker 是 host 模式来写的，因此 socket server 的地址是本机
// 需要改成指定目的 IP 的模式，在外部控制究竟与哪个 server 建立连接

// TODO
// 1. 修改成仅限目标 IP 的模式，取消对系统的判定
// 2. 研究并调试 INPUT
// 3. 加入外部对窗口大小的控制（现在是写死的 800 * 600）

#include <iostream>
#include <windows.h>
#include <vector>
#include <sstream>
#include <pthread.h>
#include <ctime>
#include <chrono>
using namespace std;

int screenWidth, screenHeight;
int server; // TODO: Move to local variable
chrono::_V2::system_clock::time_point last_ping;
bool done;
HWND hwnd;
char *winTitle;
char dockerHost[20];
string hardcodeIP;
int targetPort;

int windowHeight;
int windowWidth;

const byte MOUSE_MOVE = 0;
const byte MOUSE_DOWN = 1;
const byte MOUSE_UP = 2;
const byte KEY_UP = 0;
const byte KEY_DOWN = 1;

int clientConnect()
{
    WSADATA wsa_data;
    SOCKADDR_IN addr;

    memset(&addr, 0, sizeof(addr));
    // WSAStartup(MAKEWORD(2, 0), &wsa_data);
    WSAStartup(0x101, &wsa_data);

    addr.sin_family = AF_INET;
    addr.sin_port = htons(targetPort);

    cout << "Running with hardcode IP:" << hardcodeIP << ", port:" << targetPort << endl;
    // addr.sin_addr.s_addr = inet_addr(hardcodeIP.c_str());
    addr.sin_addr.S_un.S_addr = inet_addr(hardcodeIP.c_str());

    while(1){
        cout << "Connect to " << hardcodeIP.c_str() << ":" << targetPort << endl;
        int serverSocket = socket(AF_INET, SOCK_STREAM, 0);
        // int err = connect(serverSocket, reinterpret_cast<SOCKADDR *>(&addr), sizeof(addr));
        int err = connect(serverSocket, (struct sockaddr*)&addr, sizeof(addr));
        if ( err >= 0){
            cout << "Connection success!" << endl;
            return serverSocket;
        }
        else{
            cout << "Connection failed, try again. " << "Server socket: " << serverSocket <<" Error code: " << err << endl;
            Sleep(1000);
        }
    }
}

HWND getWindowByTitle(char *pattern)
{
    HWND hwnd = NULL;

    do
    {
        hwnd = FindWindowEx(NULL, hwnd, NULL, NULL);
        DWORD dwPID = 0;
        GetWindowThreadProcessId(hwnd, &dwPID);
        cout << "Scanning hwnd: " << hwnd << endl;
        int len = GetWindowTextLength(hwnd) + 1;

        char title[len];
        GetWindowText(hwnd, title, len);
        string st(title);
        cout << "Title: <" << st << ">" << endl;

        if (st.find(pattern) != string::npos)
        {
            cout << "Found " << hwnd << endl;
            return hwnd;
        }
    } while (hwnd != 0);
    cout << "Not found" << endl;
    return hwnd; //Ignore that
}

// isDxGame use hardware keys, it's special case
HWND sendIt(int key, bool state, bool isDxGame)
{
    cout << "Sending key " << ' ' << key << endl;
    HWND temp = SetActiveWindow(hwnd);
    ShowWindow(hwnd, SW_RESTORE);
    SetFocus(hwnd);
    BringWindowToTop(hwnd);

    INPUT ip;
    ZeroMemory(&ip, sizeof(INPUT));
    // Set up a generic keyboard event.
    ip.type = INPUT_KEYBOARD;
    ip.ki.time = 0;
    ip.ki.dwExtraInfo = 0;
    if (isDxGame)
    {
        if (key == VK_UP || key == VK_DOWN || key == VK_LEFT || key == VK_RIGHT)
        {
            ip.ki.dwFlags = KEYEVENTF_EXTENDEDKEY;
            cout << "after " << key << endl;
        }
        key = MapVirtualKey(key, 0);
        ip.ki.wScan = key; // hardware scan code for key
        ip.ki.dwFlags |= KEYEVENTF_SCANCODE;
        cout << "after " << key << endl;
    }
    else
    {
        ip.ki.wVk = key; // virtual-key code for the "a" key
    }
    if (state == KEY_UP)
    {
        ip.ki.dwFlags |= KEYEVENTF_KEYUP;
    }
    SendInput(1, &ip, sizeof(INPUT));

    cout << "sended key " << ' ' << key << endl;

    return hwnd;
}

void sendMouseDown(bool isLeft, byte state, float x, float y)
{
    double fScreenWidth = GetSystemMetrics(SM_CXSCREEN) - 1;
    double fScreenHeight = GetSystemMetrics(SM_CYSCREEN) - 1;
    if (x < 0) x = 0;
    if (y < 0) y = 0;
    if (x > int(fScreenWidth)) x = int(fScreenWidth);
    if (y > int(fScreenHeight)) y = int(fScreenHeight);
    double fx = x * (65535.0f / fScreenWidth);
    double fy = y * (65535.0f / fScreenHeight);
    INPUT Input[1] = {};
    ZeroMemory(Input, sizeof(INPUT));
    Input[0].type = INPUT_MOUSE;
    Input[0].mi.dx = long(fx);
    Input[0].mi.dy = long(fy);
    cout << "isLeft: " << isLeft << " state: " << int(state) << endl;
    cout << "dx: " << Input[0].mi.dx << " dy: " << Input[0].mi.dy << " (0 - 65535)" << endl;

    if (state == MOUSE_MOVE){
        Input[0].mi.dwFlags = MOUSEEVENTF_MOVE | MOUSEEVENTF_ABSOLUTE;
        cout << "Move." << endl;
    }
    else{
        if (isLeft && state == MOUSE_DOWN){
            Input[0].mi.dwFlags = MOUSEEVENTF_LEFTDOWN | MOUSEEVENTF_ABSOLUTE;
            cout << "Left down." << endl;
        }
        else if (isLeft && state == MOUSE_UP){
            Input[0].mi.dwFlags = MOUSEEVENTF_LEFTUP | MOUSEEVENTF_ABSOLUTE;
            cout << "Left up." << endl;
        }
        else if (!isLeft && state == MOUSE_DOWN){
            Input[0].mi.dwFlags = MOUSEEVENTF_RIGHTDOWN | MOUSEEVENTF_ABSOLUTE;
            cout << "Right down." << endl;
        }
        else if (!isLeft && state == MOUSE_UP){
            Input[0].mi.dwFlags = MOUSEEVENTF_RIGHTUP | MOUSEEVENTF_ABSOLUTE;
            cout << "Right up." << endl;
        }
    }
    
    cout << "dwFlags: " << Input[0].mi.dwFlags << endl;
    UINT uSent = SendInput(ARRAYSIZE(Input), Input, sizeof(INPUT));
    if (uSent != ARRAYSIZE(Input))
    {
        cout << "SendInput failed: 0x" << HRESULT_FROM_WIN32(GetLastError()) << endl;
    } 
}

struct Mouse
{
    byte isLeft;
    byte state;
    float x;
    float y;
    float relwidth;
    float relheight;
};

struct Key
{
    byte key;
    byte state;
};

// TODO: Use some proper serialization?
Key parseKeyPayload(string stPos)
{
    stringstream ss(stPos);     // f"{KeyCode},{KeyState}"

    string substr;
    getline(ss, substr, ',');   // f"{KeyCode}"
    byte key = stof(substr);    // int

    getline(ss, substr, ',');   // f"{KeyState}"
    byte state = stof(substr);  // bool

    return Key{key, state};
}

// TODO: Use some proper serialization?
Mouse parseMousePayload(string stPos)
{
    stringstream ss(stPos);     // f"{IsLeft},{mouseState},{X},{Y},{Width},{Height}"

    string substr;
    getline(ss, substr, ',');
    bool isLeft = stof(substr); // f"{IsLeft}"

    getline(ss, substr, ',');
    byte state = stof(substr);  // f"{mouseState}"

    getline(ss, substr, ',');
    float x = stof(substr);     // f"{X}"

    getline(ss, substr, ',');
    float y = stof(substr);     // f"{Y}"

    getline(ss, substr, ',');
    float w = stof(substr);     // f"{Width}"

    getline(ss, substr, ',');
    float h = stof(substr);     // f"{Height}"

    return Mouse{isLeft, state, x, y, w, h};
}

void formatWindow(HWND hwnd)
{
    SetWindowPos(hwnd, NULL, 0, 0, windowWidth, windowHeight, 0);
    // SetWindowLong(hwnd, GWL_STYLE, 0);
    cout << "Window formated" << endl;
}

// 检查 socket 管道是否存活
// 在 main 中接受来自 server 端的定期 ping，从而实现保活
void *thealthcheck(void *args)
{
    while (true)
    {
        cout << "Health check pipe" << endl;
        auto cur = chrono::system_clock::now();
        chrono::duration<double> elapsed_seconds = cur - last_ping;
        cout << elapsed_seconds.count() << endl;
        if (elapsed_seconds.count() > 30)
        {
            // socket is died
            cout << "Broken pipe" << endl;
            done = true;
            return NULL;
        }
        Sleep(2000);
    }
}

// 检查 hwnd 值的更新
// HWND: A handle to the window to precede the positioned window in the Z order. This parameter must be a window handle or one of the following values.
void *thwndupdate(void *args)
{
    HWND oldhwnd;
    while (true)
    {
        cout << "Finding title " << winTitle << endl;
        hwnd = getWindowByTitle(winTitle);
        if (hwnd != oldhwnd)
        {
            formatWindow(hwnd);
            cout << "Updated HWND: " << hwnd << endl;
            oldhwnd = hwnd;
        }
        Sleep(2000);
    }
}

void processEvent(string ev, bool isDxGame)
{
    if (ev[0] == 'K')           // f"K{KeyCode},{KeyState}"
    {
        Key key = parseKeyPayload(ev.substr(1, ev.length() - 1));           // input = f"{KeyCode},{KeyState}"
        cout << "Input key: " << int(key.key) << " " << int(key.state) << endl;
        sendIt(key.key, key.state, isDxGame);
    }
    else if (ev[0] == 'M')      // f"M{IsLeft},{mouseState},{X},{Y},{Width},{Height}"
    {
        Mouse mouse = parseMousePayload(ev.substr(1, ev.length() - 1));     // input = f"{IsLeft},{mouseState},{X},{Y},{Width},{Height}"
        float x = mouse.x;
        float y = mouse.y;
        cout << "Mouse moving!" << endl;
        sendMouseDown(mouse.isLeft, mouse.state, x, y);
    }
    else if (ev[0] == 'C')
    {
        // inputs used for controling
        // ...
    }
}

int main(int argc, char *argv[])
{
    // argv: vmid 应用名 是否是游戏 目标后台服务器的IP 目标后台服务器的端口 窗口宽度 窗口高度
    // winTitle = (char *)"NoneName";
    windowWidth = 800;
    windowHeight = 600;
    bool isDxGame = false;
    cout << "args" << argv << endl;
    if (argc > 1){
        winTitle = argv[1];
        cout << "Title: " << winTitle << endl;
    }
    if (argc > 2)
    {
        if (strcmp(argv[2], "game") == 0)
        {
            isDxGame = true;
            cout << "is game." << endl;
        }
    }
    if (argc > 3)
    {
        hardcodeIP = argv[3];
    }
    if (argc > 4)
    {
        targetPort = atoi(argv[4]);
    }
    if (argc > 5){
        windowWidth = atoi(argv[5]);
    }
    if (argc > 6){
        windowHeight = atoi(argv[6]);
    }

    hwnd = 0;   // hwnd=0 -> Places the window at the top of the Z order.
    cout << "width " << screenWidth << " "
         << "height " << screenHeight << endl;
    cout << "hardcode IP " << hardcodeIP << endl
         << "target Port " << targetPort << endl;

    formatWindow(hwnd);

    server = clientConnect();
    cout << "Connected " << server << endl;

    // setup socket watcher.
    // TODO: How to check if a socket pipe is broken?
    done = false;
    last_ping = chrono::system_clock::now();
    pthread_t th1;
    pthread_t th2;

    int t1 = pthread_create(&th1, NULL, thealthcheck, NULL);
    // title must be precise
    int t2 = pthread_create(&th2, NULL, thwndupdate, NULL);

    int recv_size;
    char buf[2000];
    istringstream iss("");
    string ev;

    do
    {
        if (done)
        {
            // 目前 done 仅限于 pipe 中断时发生
            // 如果 done 了，如何重启连接？
            // exit(1);
            cout << "Wrong: long time no connection!" << endl;
            done = false;
            last_ping = chrono::system_clock::now(); // 先不管保活
        }
        //Receive a reply from the server
        if ((recv_size = recv(server, buf, 1024, 0)) == SOCKET_ERROR)
        {
            puts("recv failed");
            Sleep(1000);
            continue;
        }

        char *buffer = new char[recv_size];
        memcpy(buffer, buf, recv_size);
        if (recv_size == 1)
        {
            if (buffer[0] == 0)
            {
                // Received ping
                last_ping = chrono::system_clock::now();
            };
        }

        try
        {
            stringstream ss(buffer);

            while (getline(ss, ev, '|'))
            {
                processEvent(ev, isDxGame);
            }
        }
        catch (const std::exception &e)
        {
            cout << "exception" << e.what() << endl;
        }
    } while (true);
    closesocket(server);
    cout << "Socket closed." << endl
         << endl;

    WSACleanup();

    return 0;
}
