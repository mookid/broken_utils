use std::fmt::Display;
use std::io::BufRead;
use std::io::BufReader;
use std::iter::Peekable;
use std::process;

use regex::Regex;

use termcolor::Color;
use termcolor::ColorChoice;
use termcolor::ColorSpec;
use termcolor::StandardStream;
use termcolor::WriteColor;

const USAGE: &str = r#"Usage: $BIN_NAME [options] [file...]
Sort matches in the given files.
If file is omitted or is '-', read from standard input.
Options:
  -V, --version       output version information
      --help          output help
  -r  re              sort according to the regexp match
  -f  field           sort according to the nth field (starting from 1)"#;
const BIN_NAME: &str = env!("CARGO_PKG_NAME");
const VERSION: &str = env!("CARGO_PKG_VERSION");
const EXIT_ERROR: i32 = 2;

fn usage(code: i32) -> ! {
    eprintln!("{}", USAGE.replace("$BIN_NAME", BIN_NAME));
    process::exit(code);
}

fn parsing_error(arg: impl std::fmt::Display) -> ! {
    eprintln!("parsing error: '{}'", arg);
    usage(EXIT_ERROR);
}

fn missing_arg(arg: impl std::fmt::Display) -> ! {
    eprintln!("option requires an argument: '{}'", arg);
    usage(EXIT_ERROR);
}

fn no_filter() -> ! {
    eprintln!("either re or field should be provided");
    usage(EXIT_ERROR);
}

fn invalid_field() -> ! {
    eprintln!("invalid field: indexing is 1-based");
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

enum FilterDescr {
    Re(String),
    Field(usize),
}

enum Filter {
    Re(Regex),
    Field(usize),
}

#[derive(Default)]
struct Opts {
    paths: Vec<String>,
    filter: Option<FilterDescr>,
}

fn parse_field(opts: &mut Opts, args: &mut Peekable<impl Iterator<Item = String>>) -> bool {
    let arg = args.next().unwrap();
    match args.next() {
        None => {
            missing_arg(arg);
        }
        Some(d) => {
            if let Some(d) = d.parse().ok() {
                opts.filter = Some(FilterDescr::Field(d));
            } else {
                parsing_error(arg);
            }
            true
        }
    }
}

fn parse_re(opts: &mut Opts, args: &mut Peekable<impl Iterator<Item = String>>) -> bool {
    let arg = args.next().unwrap();
    match args.next() {
        None => missing_arg(arg),
        Some(d) => {
            opts.filter = Some(FilterDescr::Re(d));
            true
        }
    }
}

fn parse_paths(opts: &mut Opts, args: &mut Peekable<impl Iterator<Item = String>>) -> bool {
    match args.next() {
        None => false,
        Some(d) => {
            opts.paths.push(d);
            true
        }
    }
}

fn parse_options(opts: &mut Opts, args: &mut Peekable<impl Iterator<Item = String>>) -> bool {
    if let Some(arg) = args.peek() {
        match &arg[..] {
            "--help" => usage(0),
            "--version" => show_version(),
            "-r" => parse_re(opts, args),
            "-f" => parse_field(opts, args),
            _ => false,
        }
    } else {
        false
    }
}

type Loc = Option<(usize, usize)>;

fn loc(from: usize, to: usize) -> Loc {
    Some((from, to))
}

type Match = (String, Loc);

fn extract_match(re: &Regex, text: &str) -> Option<Match> {
    match re.find(text) {
        None => None,
        Some(m) => Some((m.as_str().to_string(), loc(m.start(), m.end()))),
    }
}

fn extract_field(nfield: usize, text: &str) -> Option<Match> {
    if nfield == 0 {
        invalid_field();
    }
    if let Some(m) = text.split_whitespace().nth(nfield - 1) {
        Some((m.to_string(), None))
    } else {
        None
    }
}

impl Filter {
    fn extract(&self, text: &str) -> Option<Match> {
        match self {
            Filter::Re(re) => extract_match(re, text),
            Filter::Field(nfield) => extract_field(*nfield, text),
        }
    }
}

fn validate_opts(opts: &Opts) -> Filter {
    match &opts.filter {
        Some(FilterDescr::Re(re)) => {
            let re = format!("(?i){}", re);
            match Regex::new(&re) {
                Ok(ref re) => Filter::Re(re.clone()),
                Err(e) => die_error("can't compile re", e),
            }
        },
        Some(FilterDescr::Field(nfield)) => Filter::Field(*nfield),
        None => no_filter(),
    }
}

type Item = (String, String, Loc);

fn collect(f: &mut impl std::io::Read, filter: &Filter, results: &mut Vec<Item>) {
    let mut f = BufReader::new(f);
    let mut buf = String::new();
    loop {
        buf.clear();
        match f.read_line(&mut buf) {
            Ok(0) => break,
            Err(err) => die_error("read error", err),
            Ok(_) => {
                let buf = buf.trim_end().to_string();
                if let Some((m, loc)) = filter.extract(&buf) {
                    results.push((m, buf, loc))
                }
            }
        }
    }
}

fn main() {
    let mut opts = Default::default();
    let mut args = std::env::args().skip(1).peekable();
    while parse_options(&mut opts, &mut args) {}
    while parse_paths(&mut opts, &mut args) {}

    let filter = validate_opts(&opts);

    assert_eq!(args.peek(), None);
    let mut results = vec![];
    if opts.paths.is_empty() {
        opts.paths.push("-".to_string());
    }
    for path in opts.paths {
        if path == "-" {
            let stdin = std::io::stdin();
            let mut stdin = stdin.lock();
            collect(&mut stdin, &filter, &mut results);
        } else {
            match std::fs::File::open(&path) {
                Err(e) => die_error(&format!("can't open file {}", path), e),
                Ok(mut f) => collect(&mut f, &filter, &mut results),
            }
        }
    }
    results.sort();

    let mut w = StandardStream::stdout(ColorChoice::Always);
    let mut matchcolor = ColorSpec::new();
    matchcolor.set_fg(Some(Color::Red));
    matchcolor.set_bold(true);
    for (_, text, loc) in results {
        if let Some((from, to)) = loc {
            assert!(from <= to);
            assert!(to <= text.len());
            print!("{}", &text[..from]);
            let _ = w.set_color(&matchcolor);
            print!("{}", &text[from..to]);
            let _ = w.reset();
            println!("{}", &text[to..]);
        } else {
            println!("{}", text);
        }
    }
}
