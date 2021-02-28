use std::collections::hash_map::DefaultHasher;
use std::collections::HashMap;
use std::fmt::Display;
use std::hash::Hasher;
use std::io::Read;
use std::io::Write;
use std::iter::Peekable;
use std::path::PathBuf;
use std::process;

use walkdir::DirEntry;
use walkdir::WalkDir;

// TODO file filters?
const USAGE: &str = r#"Usage: $BIN_NAME [root-dirs...]
Recursively traverse directories, listing all duplicated files.
Arguments:
  root-dirs           the root directories where the traversal starts
                      if omitted, use the current working directory"#;
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

fn die_error<E: Display>(msg: &str, err: E) -> ! {
    eprintln!("{}: {}", msg, err);
    process::exit(EXIT_ERROR);
}

#[derive(Default)]
struct Opts {
    roots: Vec<String>,
}

fn add_root(opts: &mut Opts, args: &mut Peekable<impl Iterator<Item = String>>) -> bool {
    if let Some(arg) = args.next() {
        opts.roots.push(arg);
        true
    } else {
        false
    }
}

fn parse_options(opts: &mut Opts, args: &mut Peekable<impl Iterator<Item = String>>) -> bool {
    if let Some(arg) = args.peek() {
        let arg = arg.clone();
        match &*arg {
            "-h" => usage(0),
            "--help" => usage(0),
            "--version" => show_version(),
            _ => add_root(opts, args),
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

fn hash_content(filename: &PathBuf) -> Result<u64, std::io::Error> {
    let mut file = match std::fs::File::open(filename) {
        Err(err) => return Err(err),
        Ok(file) => file,
    };
    let mut hasher = DefaultHasher::new();
    let mut contents = vec![];
    if let Err(err) = file.read_to_end(&mut contents) {
        return Err(err);
    }
    hasher.write(&contents);
    Ok(hasher.finish())
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

    let stdout = std::io::stdout();
    let mut stdout = stdout.lock();

    let stderr = std::io::stderr();
    let mut stderr = stderr.lock();

    let mut roots = vec![];
    if opts.roots.is_empty() {
        roots.push(std::env::current_dir().unwrap())
    } else {
        for root_str in opts.roots {
            roots.push(root_str.into())
        }
    };
    let mut unique_files = HashMap::new();
    let mut unique_files2 = HashMap::new();

    for root in roots {
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
            let filename = entry.path().to_owned();
            unique_files.entry(m.len()).or_insert(vec![]).push(filename);
        }
    }

    let mut first = true;
    for (_, entries) in unique_files {
        if entries.len() >= 2 {
            unique_files2.clear();
            for filename in entries {
                match hash_content(&filename) {
                    Ok(hash) => {
                        unique_files2.entry(hash).or_insert(vec![]).push(filename);
                    }
                    Err(e) => ignore_broken_pipe(write!(
                        stderr,
                        "error while hashing contents for file {}: {}\n",
                        filename.display(),
                        e
                    )),
                }
            }
            for (_, entries) in &unique_files2 {
                if entries.len() >= 2 {
                    if first {
                        first = false;
                    } else {
                        ignore_broken_pipe(write!(stdout, "\n"));
                    }
                    for entry in entries {
                        ignore_broken_pipe(write!(stdout, "{}\n", entry.display()));
                    }
                }
            }
        }
    }
}
