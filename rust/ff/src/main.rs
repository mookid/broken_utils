use std::fmt::Display;
use std::io::Write;
use std::iter::Peekable;
use std::process;

use regex::Regex;

use walkdir::DirEntry;
use walkdir::WalkDir;

const USAGE: &str = r#"Usage: $BIN_NAME [[-v] filter...]
Recursively traverse directories, listing all files.
Arguments:
     filter         filter filenames matching the given regexp
  -v filter         filter out filenames matching the given regexp
  .                 strip common working directory prefix
  /                 add common working directory prefix"#;
const BIN_NAME: &str = env!("CARGO_PKG_NAME");
const VERSION: &str = env!("CARGO_PKG_VERSION");
const EXIT_ERROR: i32 = 2;

fn usage(code: i32) -> ! {
    eprintln!("{}", USAGE.replace("$BIN_NAME", BIN_NAME));
    process::exit(code);
}

fn missing_arg(arg: impl std::fmt::Display) -> ! {
    eprintln!("option requires an argument: '{}'", arg);
    usage(EXIT_ERROR);
}

fn show_version() -> ! {
    eprintln!("{} {}", BIN_NAME, VERSION);
    process::exit(0);
}

fn die_error<E: Display>(msg: &str, err: E) -> ! {
    eprintln!("{}: {}", msg, err);
    process::exit(EXIT_ERROR);
}

#[derive(Default)]
struct Opts {
    filters: Vec<Regex>,
    filters_neg: Vec<Regex>,
    strip_prefix: bool,
}

fn compile_re(arg: &str) -> Regex {
    match Regex::new(arg) {
        Ok(re) => re,
        Err(e) => die_error("can't compile re", e),
    }
}

fn add_filter(opts: &mut Opts, args: &mut Peekable<impl Iterator<Item = String>>) -> bool {
    if let Some(arg) = args.next() {
        opts.filters.push(compile_re(&arg));
        true
    } else {
        false
    }
}

fn strip_prefix(opts: &mut Opts, args: &mut Peekable<impl Iterator<Item = String>>) -> bool {
    if let Some(arg) = args.next() {
        opts.strip_prefix = &arg == ".";
        true
    } else {
        false
    }
}

fn add_filter_neg(opts: &mut Opts, args: &mut Peekable<impl Iterator<Item = String>>) -> bool {
    let arg = args.next().unwrap();
    match args.next() {
        Some(d) => {
            opts.filters_neg.push(compile_re(&d));
            true
        }
        None => missing_arg(arg),
    }
}

fn parse_options(opts: &mut Opts, args: &mut Peekable<impl Iterator<Item = String>>) -> bool {
    if let Some(arg) = args.peek() {
        let arg = arg.clone();
        match &arg[..] {
            "-h" => usage(0),
            "--help" => usage(0),
            "--version" => show_version(),
            "." | "/" => strip_prefix(opts, args),
            "-v" => add_filter_neg(opts, args),
            _ => add_filter(opts, args),
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

fn main() {
    let mut opts = Default::default();
    let mut args = std::env::args().skip(1).peekable();
    while parse_options(&mut opts, &mut args) {}

    let root = std::env::current_dir().unwrap();
    let stdout = std::io::stdout();
    let mut stdout = stdout.lock();

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
        let filename = entry.path();
        let filename = if opts.strip_prefix {
            filename.strip_prefix(&root).ok().unwrap_or(filename)
        } else {
            filename
        };
        let filename = filename.display().to_string();
        let filename = filename.replace("\\", "/");

        if opts.filters.iter().all(|f| f.is_match(&filename))
            && !opts.filters_neg.iter().any(|f| f.is_match(&filename))
        {
            if let Err(err) = writeln!(stdout, "{}", filename) {
                if err.kind() == std::io::ErrorKind::BrokenPipe {
                    return;
                }
                die_error("write", err)
            }
        }
    }
}
