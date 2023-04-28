package shmipc

//func TestEpoll(t *testing.T) {
//	// Create an epoll instance
//	epfd, err := syscall.EpollCreate(0)
//	if err != nil {
//		t.Fatalf("failed to create epoll instance: %v", err)
//	}
//	defer syscall.Close(epfd)
//
//	// Create a [socket pair](poe://www.poe.com/_api/key_phrase?phrase=socket%20pair&prompt=Tell%20me%20more%20about%20socket%20pair.)
//	sockpair, err := syscall.Socketpair(syscall.AF_UNIX, syscall.SOCK_STREAM, 0)
//	if err != nil {
//		t.Fatalf("failed to create socket pair: %v", err)
//	}
//	defer func() {
//		syscall.Close(sockpair[0])
//		syscall.Close(sockpair[1])
//	}()
//
//	// Add the read end of the socket pair to the epoll instance
//	var ev epollEvent
//	ev.events = syscall.EPOLLIN | epollModeET
//	ev.data[0], ev.data[1], ev.data[2], ev.data[3], ev.data[4], ev.data[5], ev.data[6], ev.data[7] = byte(sockpair[0]), 0, 0, 0, 0, 0, 0, 0
//	if err := epollCtl(epfd, syscall.EPOLL_CTL_ADD, sockpair[0], &ev); err != nil {
//		t.Fatalf("failed to add socket to epoll instance: %v", err)
//	}
//
//	// Write data to the write end of the socket pair
//	msg := []byte("hello")
//	if _, err := syscall.Write(sockpair[1], msg); err != nil {
//		t.Fatalf("failed to write data to socket: %v", err)
//	}
//
//	// Wait for the data to become available for reading
//	events := make([]epollEvent, 1)
//	if n, err := epollWait(epfd, events, -1); err != nil {
//		t.Fatalf("failed to wait for events: %v", err)
//	} else if n != 1 {
//		t.Fatalf("unexpected number of events: got %d, want 1", n)
//	}
//
//	// Read the data from the socket
//	buf := make([]byte, 5)
//	if _, err := syscall.Read(sockpair[0], buf); err != nil {
//		t.Fatalf("failed to read data from socket: %v", err)
//	} else if string(buf) != "hello" {
//		t.Fatalf("unexpected message received: got %q, want %q", string(buf), "hello")
//	}
//}
