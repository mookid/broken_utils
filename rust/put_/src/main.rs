use std::fmt::Display;
use std::io::Write;
use std::iter::Peekable;
use std::process;

use walkdir::DirEntry;
use walkdir::WalkDir;

const USAGE: &str = r#"Usage: $BIN_NAME [root-dirs...]
Recursively traverse directories, renaming all whitespace in filenames with underscores."#;
const BIN_NAME: &str = env!("CARGO_PKG_NAME");
const VERSION: &str = env!("CARGO_PKG_VERSION");
const EXIT_ERROR: i32 = 2;

fn usage(code: i32) -> ! {
    eprintln!("{}", USAGE.replace("$BIN_NAME", BIN_NAME));
    process::exit(code);
}

fn show_version() -> ! {
    eprintln!("{} {}", BIN_NAME, VERSION);
    process::exit(0);
}

fn invalid_opt(arg: impl std::fmt::Display) -> ! {
    eprintln!("invalid option: '{}'", arg);
    usage(EXIT_ERROR);
}

fn die_error<E: Display>(msg: &str, err: E) -> ! {
    eprintln!("{}: {}", msg, err);
    process::exit(EXIT_ERROR);
}

#[derive(Default)]
struct Opts {}

fn parse_options(_opts: &mut Opts, args: &mut Peekable<impl Iterator<Item = String>>) -> bool {
    if let Some(arg) = args.peek() {
        match &arg[..] {
            "-h" => usage(0),
            "--help" => usage(0),
            "--version" => show_version(),
            arg => invalid_opt(arg),
        }
    } else {
        false
    }
}

fn is_hidden(entry: &DirEntry) -> bool {
    entry
        .file_name()
        .to_str()
        .map(|s| s.starts_with("."))
        .unwrap_or(false)
}

fn ignore_broken_pipe<T>(r: std::io::Result<T>) {
    if let Err(err) = r {
        if err.kind() == std::io::ErrorKind::BrokenPipe {
            return;
        }
        die_error("write", err)
    }
}

fn main() {
    let mut opts = Default::default();
    let mut args = std::env::args().skip(1).peekable();
    while parse_options(&mut opts, &mut args) {}

    let stderr = std::io::stderr();
    let mut stderr = stderr.lock();

    let root = std::env::current_dir().unwrap();

    let entries = WalkDir::new(&root).into_iter();
    for entry in entries.filter_entry(|e| !is_hidden(e)) {
        let entry = match entry {
            Ok(e) => e,
            Err(_) => continue,
        };
        let m = match entry.metadata() {
            Ok(m) => m,
            Err(_) => continue,
        };
        if m.is_dir() {
            continue;
        }
        let filename = if let Some(f) = entry.path().to_str() {
            f
        } else {
            continue;
        };
        if filename.as_bytes().iter().any(|c| *c == b' ') {
            let new_name = filename.replace(" ", "_");
            let w = if let Err(err) = std::fs::rename(filename, &new_name) {
                write!(
                    stderr,
                    "ERROR: failed to rename file {}: {}\n",
                    filename, err
                )
            } else {
                write!(stderr, "renamed file {} -> {}\n", filename, new_name)
            };
            ignore_broken_pipe(w);
        }
    }
}
