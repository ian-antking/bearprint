import textwrap

class Printer:
    def __init__(self, device='/dev/usb/lp0', encoding='utf-8'):
        self.device = device
        self.encoding = encoding
        self.line_width = 32

    def _write(self, data: bytes):
        with open(self.device, 'wb') as printer:
            printer.write(data)

    def text(self, message: str):
      wrapped_lines = textwrap.wrap(message, self.line_width)
      for line in wrapped_lines:
          self._write(line.encode(self.encoding) + b'\n')

    def blank_line(self, count=1):
        self._write(b'\n' * count)

    def cut(self):
        self._write(b'\x1D\x56\x00')

    def print_text(self, message: str, blank_lines=6, cut_paper=True):
        self.text(message)
        self.blank_line(blank_lines)
        if cut_paper:
            self.cut()
