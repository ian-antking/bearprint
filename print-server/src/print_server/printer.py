import textwrap

class Printer:
    def __init__(self, device='/dev/usb/lp0', encoding='cp437'):
        self.device = device
        self.encoding = encoding
        self.line_width = 64

    def _write(self, data: bytes):
        with open(self.device, 'wb') as printer:
            printer.write(data)

    def _write_line(self, line: str):
        self._write(line.encode(self.encoding, errors='replace') + b'\n')

    def format_line(self, line: str, align: str = "left") -> str:
        if align == "center":
            return line.center(self.line_width)
        elif align == "right":
            return line.rjust(self.line_width)
        else:
            return line.ljust(self.line_width)

    def text(self, message: str, align: str = "left"):
        for raw_line in message.split('\n'):
            wrapped_lines = textwrap.wrap(raw_line, self.line_width)
            for line in wrapped_lines or ['']:
                self._write_line(self.format_line(line, align))

    def blank_line(self, count=1):
        self._write(b'\n' * count)

    def cut(self):
        self.blank_line(6)
        self._write(b'\x1D\x56\x00')

    def print_job(self, job):
        for item in job:
            type_ = item.get("type")
            if type_ == "text":
                self.text(item.get("content", ""), align=item.get("align", "left"))
            elif type_ == "blank":
                self.blank_line(item.get("count", 1))
            elif type_ == "line":
                self.text("-" * self.line_width)
            elif type_ == "cut":
                self.cut()
