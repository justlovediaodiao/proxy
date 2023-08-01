import signal
import subprocess
import sys

from PyQt6.QtWidgets import (QApplication, QComboBox, QLabel, QVBoxLayout,
                             QWidget)


class MainWindow(QWidget):

    def __init__(self):
        super().__init__()
        self.proxy = Proxy()

        self.combo_box = QComboBox(self)
        self.combo_box.addItem("Off")
        self.combo_box.addItem("Global")
        self.combo_box.addItem("Pac")
        self.combo_box.currentIndexChanged.connect(self.on_combo_box_changed)

        self.label = QLabel(self)

        layout = QVBoxLayout()
        layout.addWidget(self.combo_box)
        layout.addWidget(self.label)
        self.setLayout(layout)

        self.setWindowTitle("Proxy")
    
    def success(self, msg):
        self.label.setText(msg)
        self.label.setStyleSheet("color: green")
    
    def error(self, msg):
        self.label.setText(msg)
        self.label.setStyleSheet("color: red")

    def on_combo_box_changed(self, index):
        text = self.combo_box.currentText()
        try:
            match text:
                case "Off":
                    self.proxy.off()
                    self.success("Proxy off")
                case "Global":
                    self.proxy.global_()
                    self.success("Global mode")
                case "Pac":
                    self.proxy.pac()
                    self.success("Pac mod")
        except FileNotFoundError:
            self.error("Error: proxy not found")
        except Exception as e:
            self.error(f"Error: {e}")


class Proxy:

    def __init__(self):
        self._process = None

    def _start(self, cmd: list):
        self._process = subprocess.Popen(cmd)

    def _stop(self):
        if self._process:
            self._process.send_signal(signal.SIGINT)
            self._process.wait()
            self._process = None

    def off(self):
        self._stop()
        subprocess.run(["./proxy", "off"])

    def global_(self):
        self._stop()
        self._start(["./proxy", "global"])

    def pac(self):
        self._stop()
        self._start(["./proxy", "global"])


if __name__ == "__main__":
    app = QApplication(sys.argv)
    window = MainWindow()
    window.show()
    app.aboutToQuit.connect(window.proxy.off)
    sys.exit(app.exec())
