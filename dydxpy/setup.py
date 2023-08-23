from setuptools import find_namespace_packages, setup

with open('requirements.txt') as f:
    required = f.read().splitlines()

setup(
    name="dydxpy",
    version="{{VERSION}}",
    author="dYdX Trading Inc.",
    author_email="contact@dydx.exchange",
    description="Protos for dYdX v4 protocol",
    packages = find_namespace_packages(),
    install_requires=required,
    python_requires=">=3.8",
)
