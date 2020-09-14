use std::iter::Peekable;
use std::process;
use std::time::Duration;

const USAGE: &str = r#"Usage: $BIN_NAME [options] command [arg...]
Options:
  -V, --version       output version information
      --help          output help
  -b, --beep          interrupt if the program exits with non-zero code
  -n, --interval ms   set interval duration (default: 2000 ms)"#;
const BIN_NAME: &str = env!("CARGO_PKG_NAME");
const VERSION: &str = env!("CARGO_PKG_VERSION");
const EXIT_ERROR: i32 = 2;

fn ignore<T>(_: T) {}

fn usage(code: i32) -> ! {
    eprintln!("{}", USAGE.replace("$BIN_NAME", BIN_NAME));
    process::exit(code);
}

fn parsing_error(arg: impl std::fmt::Display) -> ! {
    eprintln!("parsing error: '{}'", arg);
    usage(EXIT_ERROR);
}

fn invalid_opt(arg: impl std::fmt::Display) -> ! {
    eprintln!("invalid option: '{}'", arg);
    usage(EXIT_ERROR);
}

fn missing_arg(arg: impl std::fmt::Display) -> ! {
    eprintln!("option requires an argument: '{}'", arg);
    usage(EXIT_ERROR);
}

fn show_version() -> ! {
    eprintln!("{} {}", BIN_NAME, VERSION);
    process::exit(0);
}

fn die_io_error(msg: &'static str, e: std::io::Error) -> ! {
    eprintln!("{}: {}", msg, e);
    process::exit(EXIT_ERROR);
}

#[derive(Default)]
struct Opts {
    interval: Option<u64>,
    beeps: bool,
}

fn parse_beeps(opts: &mut Opts, args: &mut Peekable<impl Iterator<Item = String>>) -> bool {
    args.next();
    opts.beeps = true;
    true
}

fn parse_interval(opts: &mut Opts, args: &mut Peekable<impl Iterator<Item = String>>) -> bool {
    let arg = args.next().unwrap();
    match args.next() {
        None => {
            missing_arg(arg);
        }
        Some(d) => {
            if let Some(d) = d.parse().ok() {
                opts.interval = Some(d);
            } else {
                parsing_error(arg);
            }
            true
        }
    }
}

fn parse_arg(opts: &mut Opts, args: &mut Peekable<impl Iterator<Item = String>>) -> bool {
    if let Some(arg) = args.peek() {
        let arg = arg.clone();
        match &*arg {
            "--help" => usage(0),
            "--version" => show_version(),
            "--beeps" => parse_beeps(opts, args),
            "--interval" => parse_interval(opts, args),
            _ if arg.starts_with("--") => invalid_opt(arg),
            _ if !arg.starts_with("-") => false,
            "-n" => parse_interval(opts, args),
            _ => {
                for ch in arg.chars().skip(1) {
                    match ch {
                        'V' => show_version(),
                        'b' => ignore(parse_beeps(opts, args)),
                        ch => invalid_opt(ch),
                    }
                }
                true
            }
        }
    } else {
        false
    }
}

fn main() {
    let mut opts = Default::default();
    let mut args = std::env::args().skip(1).peekable();
    while parse_arg(&mut opts, &mut args) {}

    let mut p = match args.next() {
        None => usage(EXIT_ERROR),
        Some(pname) => process::Command::new(pname),
    };
    for arg in args {
        p.arg(arg);
    }

    let interval = opts.interval.unwrap_or(2000);

    loop {
        if let Err(e) = p.spawn() {
            die_io_error("failed to spawn process", e);
        }
        match p.output() {
            Ok(out) if opts.beeps && !out.status.success() => process::exit(EXIT_ERROR),
            Ok(_) => {}
            Err(e) => die_io_error("error waiting for process", e),
        }
        std::thread::sleep(Duration::from_millis(interval));
    }
}
