CARGO_BIN=$(HOME)/.cargo/bin

all:
	cargo build --release

install:
	cp target/release/cmdtime.exe $(CARGO_BIN)/cmdtime.exe
	cp target/release/re-sort.exe $(CARGO_BIN)/re-sort.exe
	cp target/release/repeat.exe $(CARGO_BIN)/repeat.exe
	cp target/release/trimcolor.exe $(CARGO_BIN)/trimcolor.exe

clean:
	cargo clean
