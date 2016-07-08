package colfer

import (
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// GenerateJava writes the code into the respective ".java" files.
func GenerateJava(basedir string, packages []*Package) error {
	t := template.New("java-code").Delims("<:", ":>")
	template.Must(t.Parse(javaCode))

	for _, p := range packages {
		p.NameNative = strings.Replace(p.Name, "/", ".", -1)
	}

	for _, p := range packages {
		pkgdir, err := makePkgDir(p, basedir)
		if err != nil {
			return err
		}

		for _, s := range p.Structs {
			for _, f := range s.Fields {
				switch f.Type {
				default:
					if f.TypeRef == nil {
						f.TypeNative = f.Type
					} else {
						f.TypeNative = f.TypeRef.NameTitle()
						if f.TypeRef.Pkg != p {
							f.TypeNative = f.TypeRef.Pkg.NameNative + "." + f.TypeNative
						}
					}
				case "bool":
					f.TypeNative = "boolean"
				case "uint32", "int32":
					f.TypeNative = "int"
				case "uint64", "int64":
					f.TypeNative = "long"
				case "float32":
					f.TypeNative = "float"
				case "float64":
					f.TypeNative = "double"
				case "timestamp":
					f.TypeNative = "java.time.Instant"
				case "text":
					f.TypeNative = "String"
				case "binary":
					f.TypeNative = "byte[]"
				}
			}

			f, err := os.Create(filepath.Join(pkgdir, s.NameTitle()+".java"))
			if err != nil {
				return err
			}
			defer f.Close()

			if err := t.Execute(f, s); err != nil {
				return err
			}
		}
	}
	return nil
}

const javaCode = `package <:.Pkg.NameNative:>;


// This file was generated by colf(1); DO NOT EDIT


import static java.lang.String.format;
import java.util.InputMismatchException;
import java.nio.BufferOverflowException;
import java.nio.BufferUnderflowException;


/**
 * Data bean with built-in serialization support.
 * @author generated by colf(1)
 * @see <a href="https://github.com/pascaldekloe/colfer">Colfer's home</a>
 */
public class <:.NameTitle:> implements java.io.Serializable {

	/** The upper limit for serial byte sizes. */
	public static int colferSizeMax = 16 * 1024 * 1024;

	/** The upper limit for the number of elements in a list. */
	public static int colferListMax = 64 * 1024;

	private static final java.nio.charset.Charset _utf8 = java.nio.charset.Charset.forName("UTF-8");
<:- range .Fields:>
<:- if eq .Type "binary":>
	private static final byte[] _zero<:.NameTitle:> = new byte[0];
<:- else if .TypeArray:>
	private static final <:.TypeNative:>[] _zero<:.NameTitle:> = new <:.TypeNative:>[0];
<:- end:>
<:- end:>
<:range .Fields:>
	public <:.TypeNative:><:if .TypeArray:>[]<:end:> <:.Name:>
<:- if eq .Type "text":> = ""
<:- else if eq .Type "binary":> = _zero<:.NameTitle:>
<:- else if .TypeArray:> = _zero<:.NameTitle:>
<:- end:>;<:end:>


	/**
	 * Serializes the object.
<:- range .Fields:><:if .TypeArray:>
	 * All {@code null} entries in {@link #<:.Name:>} will be replaced with a {@code new} value.
<:- end:><:end:>
	 * @param buf the data destination.
	 * @param offset the initial index for {@code buf}, inclusive.
	 * @return the final index for {@code buf}, exclusive.
	 * @throws BufferOverflowException when {@code buf} is too small.
	 * @throws IllegalStateException on an upper limit breach defined by either {@link #colferSizeMax} or {@link #colferListMax}.
	 */
	public int marshal(byte[] buf, int offset) {
		int i = offset;
		try {
<:- range .Fields:><:if eq .Type "bool":>
			if (this.<:.Name:>) {
				buf[i++] = (byte) <:.Index:>;
			}
<:else if eq .Type "uint32":>
			if (this.<:.Name:> != 0) {
				int x = this.<:.Name:>;
				if ((x & ~((1 << 21) - 1)) != 0) {
					buf[i++] = (byte) (<:.Index:> | 0x80);
					buf[i++] = (byte) (x >>> 24);
					buf[i++] = (byte) (x >>> 16);
					buf[i++] = (byte) (x >>> 8);
					buf[i++] = (byte) (x);
				} else {
					buf[i++] = (byte) <:.Index:>;
					while (x > 0x7f) {
						buf[i++] = (byte) (x | 0x80);
						x >>>= 7;
					}
					buf[i++] = (byte) x;
				}
			}
<:else if eq .Type "uint64":>
			if (this.<:.Name:> != 0) {
				long x = this.<:.Name:>;
				if ((x & ~((1 << 49) - 1)) != 0) {
					buf[i++] = (byte) (<:.Index:> | 0x80);
					buf[i++] = (byte) (x >>> 56);
					buf[i++] = (byte) (x >>> 48);
					buf[i++] = (byte) (x >>> 40);
					buf[i++] = (byte) (x >>> 32);
					buf[i++] = (byte) (x >>> 24);
					buf[i++] = (byte) (x >>> 16);
					buf[i++] = (byte) (x >>> 8);
					buf[i++] = (byte) (x);
				} else {
					buf[i++] = (byte) <:.Index:>;
					while (x > 0x7fL) {
						buf[i++] = (byte) (x | 0x80);
						x >>>= 7;
					}
					buf[i++] = (byte) x;
				}
			}
<:else if eq .Type "int32":>
			if (this.<:.Name:> != 0) {
				int x = this.<:.Name:>;
				if (x < 0) {
					x = -x;
					buf[i++] = (byte) (<:.Index:> | 0x80);
				} else
					buf[i++] = (byte) <:.Index:>;
				while ((x & ~0x7f) != 0) {
					buf[i++] = (byte) (x | 0x80);
					x >>>= 7;
				}
				buf[i++] = (byte) x;
			}
<:else if eq .Type "int64":>
			if (this.<:.Name:> != 0) {
				long x = this.<:.Name:>;
				if (x < 0) {
					x = -x;
					buf[i++] = (byte) (<:.Index:> | 0x80);
				} else
					buf[i++] = (byte) <:.Index:>;
				for (int n = 0; n < 8 && (x & ~0x7fL) != 0; n++) {
					buf[i++] = (byte) (x | 0x80);
					x >>>= 7;
				}
				buf[i++] = (byte) x;
			}
<:else if eq .Type "float32":>
			if (this.<:.Name:> != 0.0f) {
				buf[i++] = (byte) <:.Index:>;
				int x = Float.floatToRawIntBits(this.<:.Name:>);
				buf[i++] = (byte) (x >>> 24);
				buf[i++] = (byte) (x >>> 16);
				buf[i++] = (byte) (x >>> 8);
				buf[i++] = (byte) (x);
			}
<:else if eq .Type "float64":>
			if (this.<:.Name:> != 0.0) {
				buf[i++] = (byte) <:.Index:>;
				long x = Double.doubleToRawLongBits(this.<:.Name:>);
				buf[i++] = (byte) (x >>> 56);
				buf[i++] = (byte) (x >>> 48);
				buf[i++] = (byte) (x >>> 40);
				buf[i++] = (byte) (x >>> 32);
				buf[i++] = (byte) (x >>> 24);
				buf[i++] = (byte) (x >>> 16);
				buf[i++] = (byte) (x >>> 8);
				buf[i++] = (byte) (x);
			}
<:else if eq .Type "timestamp":>
			if (this.<:.Name:> != null) {
				long s = this.<:.Name:>.getEpochSecond();
				int ns = this.<:.Name:>.getNano();
				if (s != 0 || ns != 0) {
					if (s >= 0 && s < (1L << 32)) {
						buf[i++] = (byte) <:.Index:>;
						buf[i++] = (byte) (s >>> 24);
						buf[i++] = (byte) (s >>> 16);
						buf[i++] = (byte) (s >>> 8);
						buf[i++] = (byte) (s);
						buf[i++] = (byte) (ns >>> 24);
						buf[i++] = (byte) (ns >>> 16);
						buf[i++] = (byte) (ns >>> 8);
						buf[i++] = (byte) (ns);
					} else {
						buf[i++] = (byte) (<:.Index:> | 0x80);
						buf[i++] = (byte) (s >>> 56);
						buf[i++] = (byte) (s >>> 48);
						buf[i++] = (byte) (s >>> 40);
						buf[i++] = (byte) (s >>> 32);
						buf[i++] = (byte) (s >>> 24);
						buf[i++] = (byte) (s >>> 16);
						buf[i++] = (byte) (s >>> 8);
						buf[i++] = (byte) (s);
						buf[i++] = (byte) (ns >>> 24);
						buf[i++] = (byte) (ns >>> 16);
						buf[i++] = (byte) (ns >>> 8);
						buf[i++] = (byte) (ns);
					}
				}
			}
<:else if eq .Type "text":>
			if (! this.<:.Name:>.isEmpty()) {
				buf[i++] = (byte) <:.Index:>;
				int start = ++i;

				String s = this.<:.Name:>;
				for (int sIndex = 0, sLength = s.length(); sIndex < sLength; sIndex++) {
					char c = s.charAt(sIndex);
					if (c < '\u0080') {
						buf[i++] = (byte) c;
					} else if (c < '\u0800') {
						buf[i++] = (byte) (192 | c >>> 6);
						buf[i++] = (byte) (128 | c & 63);
					} else if (c < '\ud800' || c > '\udfff') {
						buf[i++] = (byte) (224 | c >>> 12);
						buf[i++] = (byte) (128 | c >>> 6 & 63);
						buf[i++] = (byte) (128 | c & 63);
					} else {
						int cp = 0;
						if (++sIndex < sLength) cp = Character.toCodePoint(c, s.charAt(sIndex));
						if ((cp >= 1 << 16) && (cp < 1 << 21)) {
							buf[i++] = (byte) (240 | cp >>> 18);
							buf[i++] = (byte) (128 | cp >>> 12 & 63);
							buf[i++] = (byte) (128 | cp >>> 6 & 63);
							buf[i++] = (byte) (128 | cp & 63);
						} else
							buf[i++] = (byte) '?';
					}
				}
				int size = i - start;
				if (size > colferSizeMax)
					throw new IllegalStateException(format("colfer: field <:.String:> size %d exceeds %d UTF-8 bytes", size, colferSizeMax));

				int ii = start - 1;
				if (size > 0x7f) {
					i++;
					for (int x = size; x >= 1 << 14; x >>>= 7) i++;
					System.arraycopy(buf, start, buf, i - size, size);

					do {
						buf[ii++] = (byte) (size | 0x80);
						size >>>= 7;
					} while (size > 0x7f);
				}
				buf[ii] = (byte) size;
			}
<:else if eq .Type "binary":>
			if (this.<:.Name:>.length != 0) {
				buf[i++] = (byte) <:.Index:>;

				int size = this.<:.Name:>.length;
				if (size > colferSizeMax)
					throw new IllegalStateException(format("colfer: field <:.String:> size %d exceeds %d bytes", size, colferSizeMax));

				int x = size;
				while (x > 0x7f) {
					buf[i++] = (byte) (x | 0x80);
					x >>>= 7;
				}
				buf[i++] = (byte) x;

				int start = i;
				i += size;
				System.arraycopy(this.<:.Name:>, 0, buf, start, size);
			}
<:else if .TypeArray:>
			if (this.<:.Name:>.length != 0) {
				buf[i++] = (byte) <:.Index:>;
				<:.TypeNative:>[] a = this.<:.Name:>;

				int x = a.length;
				if (x > colferListMax)
					throw new IllegalStateException(format("colfer: field <:.String:> length %d exceeds %d elements", x, colferListMax));
				while (x > 0x7f) {
					buf[i++] = (byte) (x | 0x80);
					x >>>= 7;
				}
				buf[i++] = (byte) x;

				for (int ai = 0; ai < a.length; ai++) {
					<:.TypeNative:> o = a[ai];
					if (o == null) {
						o = new <:.TypeNative:>();
						a[ai] = o;
					}
					i = o.marshal(buf, i);
				}
			}
<:else:>
			if (this.<:.Name:> != null) {
				buf[i++] = (byte) <:.Index:>;
				i = this.<:.Name:>.marshal(buf, i);
			}
<:end:><:end:>
			buf[i++] = (byte) 0x7f;
			return i;
		} catch (IndexOutOfBoundsException e) {
			if (i - offset > colferSizeMax)
				throw new IllegalStateException(format("colfer: serial exceeds %d bytes", colferSizeMax));
			if (i >= buf.length)
				throw new BufferOverflowException();
			throw e;
		}
	}

	/**
	 * Deserializes the object.
	 * @param buf the data source.
	 * @param offset the initial index for {@code buf}, inclusive.
	 * @return the final index for {@code buf}, exclusive.
	 * @throws BufferUnderflowException when {@code buf} is incomplete. (EOF)
	 * @throws SecurityException on an upper limit breach defined by either {@link #colferSizeMax} or {@link #colferListMax}.
	 * @throws InputMismatchException when the data does not match this object's schema.
	 */
	public int unmarshal(byte[] buf, int offset) {
		int i = offset;
		try {
			byte header = buf[i++];
<:range .Fields:><:if eq .Type "bool":>
			if (header == (byte) <:.Index:>) {
				this.<:.Name:> = true;
				header = buf[i++];
			}
<:else if eq .Type "uint32":>
			if (header == (byte) <:.Index:>) {
				int x = 0;
				for (int shift = 0; true; shift += 7) {
					byte b = buf[i++];
					x |= (b & 0x7f) << shift;
					if (shift == 28 || b >= 0) break;
				}
				this.<:.Name:> = x;
				header = buf[i++];
			} else if (header == (byte) (<:.Index:> | 0x80)) {
				this.<:.Name:> = (buf[i++] & 0xff) << 24 | (buf[i++] & 0xff) << 16 | (buf[i++] & 0xff) << 8 | (buf[i++] & 0xff);
				header = buf[i++];
			}
<:else if eq .Type "uint64":>
			if (header == (byte) <:.Index:>) {
				long x = 0;
				for (int shift = 0; true; shift += 7) {
					byte b = buf[i++];
					if (shift == 56 || b >= 0) {
						x |= (b & 0xffL) << shift;
						break;
					}
					x |= (b & 0x7fL) << shift;
				}
				this.<:.Name:> = x;
				header = buf[i++];
			} else if (header == (byte) (<:.Index:> | 0x80)) {
				this.<:.Name:> = (buf[i++] & 0xffL) << 56 | (buf[i++] & 0xffL) << 48 | (buf[i++] & 0xffL) << 40 | (buf[i++] & 0xffL) << 32
					| (buf[i++] & 0xffL) << 24 | (buf[i++] & 0xffL) << 16 | (buf[i++] & 0xffL) << 8 | (buf[i++] & 0xffL);
				header = buf[i++];
			}
<:else if eq .Type "int32":>
			if (header == (byte) <:.Index:>) {
				int x = 0;
				for (int shift = 0; true; shift += 7) {
					byte b = buf[i++];
					x |= (b & 0x7f) << shift;
					if (shift == 28 || b >= 0) break;
				}
				this.<:.Name:> = x;
				header = buf[i++];
			} else if (header == (byte) (<:.Index:> | 0x80)) {
				int x = 0;
				for (int shift = 0; true; shift += 7) {
					byte b = buf[i++];
					x |= (b & 0x7f) << shift;
					if (shift == 28 || b >= 0) break;
				}
				this.<:.Name:> = -x;
				header = buf[i++];
			}
<:else if eq .Type "int64":>
			if (header == (byte) <:.Index:>) {
				long x = 0;
				for (int shift = 0; true; shift += 7) {
					byte b = buf[i++];
					if (shift == 56 || b >= 0) {
						x |= (b & 0xffL) << shift;
						break;
					}
					x |= (b & 0x7fL) << shift;
				}
				this.<:.Name:> = x;
				header = buf[i++];
			} else if (header == (byte) (<:.Index:> | 0x80)) {
				long x = 0;
				for (int shift = 0; true; shift += 7) {
					byte b = buf[i++];
					if (shift == 56 || b >= 0) {
						x |= (b & 0xffL) << shift;
						break;
					}
					x |= (b & 0x7fL) << shift;
				}
				this.<:.Name:> = -x;
				header = buf[i++];
			}
<:else if eq .Type "float32":>
			if (header == (byte) <:.Index:>) {
				int x = (buf[i++] & 0xff) << 24 | (buf[i++] & 0xff) << 16 | (buf[i++] & 0xff) << 8 | (buf[i++] & 0xff);
				this.<:.Name:> = Float.intBitsToFloat(x);
				header = buf[i++];
			}
<:else if eq .Type "float64":>
			if (header == (byte) <:.Index:>) {
				long x = (buf[i++] & 0xffL) << 56 | (buf[i++] & 0xffL) << 48 | (buf[i++] & 0xffL) << 40 | (buf[i++] & 0xffL) << 32
					| (buf[i++] & 0xffL) << 24 | (buf[i++] & 0xffL) << 16 | (buf[i++] & 0xffL) << 8 | (buf[i++] & 0xffL);
				this.<:.Name:> = Double.longBitsToDouble(x);
				header = buf[i++];
			}
<:else if eq .Type "timestamp":>
			if (header == (byte) <:.Index:>) {
				long s = (buf[i++] & 0xffL) << 24 | (buf[i++] & 0xffL) << 16 | (buf[i++] & 0xffL) << 8 | (buf[i++] & 0xffL);
				long ns = (buf[i++] & 0xffL) << 24 | (buf[i++] & 0xffL) << 16 | (buf[i++] & 0xffL) << 8 | (buf[i++] & 0xffL);
				this.<:.Name:> = java.time.Instant.ofEpochSecond(s, ns);
				header = buf[i++];
			} else if (header == (byte) (<:.Index:> | 0x80)) {
				long s = (buf[i++] & 0xffL) << 56 | (buf[i++] & 0xffL) << 48 | (buf[i++] & 0xffL) << 40 | (buf[i++] & 0xffL) << 32
					| (buf[i++] & 0xffL) << 24 | (buf[i++] & 0xffL) << 16 | (buf[i++] & 0xffL) << 8 | (buf[i++] & 0xffL);
				long ns = (buf[i++] & 0xffL) << 24 | (buf[i++] & 0xffL) << 16 | (buf[i++] & 0xffL) << 8 | (buf[i++] & 0xffL);
				this.<:.Name:> = java.time.Instant.ofEpochSecond(s, ns);
				header = buf[i++];
			}
<:else if eq .Type "text":>
			if (header == (byte) <:.Index:>) {
				int size = 0;
				for (int shift = 0; true; shift += 7) {
					byte b = buf[i++];
					size |= (b & 0x7f) << shift;
					if (shift == 28 || b >= 0) break;
				}
				if (size > colferSizeMax)
					throw new SecurityException(format("colfer: field <:.String:> size %d exceeds %d UTF-8 bytes", size, colferSizeMax));

				int start = i;
				i += size;
				this.<:.Name:> = new String(buf, start, size, this._utf8);
				header = buf[i++];
			}
<:else if eq .Type "binary":>
			if (header == (byte) <:.Index:>) {
				int size = 0;
				for (int shift = 0; true; shift += 7) {
					byte b = buf[i++];
					size |= (b & 0x7f) << shift;
					if (shift == 28 || b >= 0) break;
				}
				if (size > colferSizeMax)
					throw new SecurityException(format("colfer: field <:.String:> size %d exceeds %d bytes", size, colferSizeMax));

				this.<:.Name:> = new byte[size];
				int start = i;
				i += size;
				System.arraycopy(buf, start, this.<:.Name:>, 0, size);
				header = buf[i++];
			}
<:else if .TypeArray:>
			if (header == (byte) <:.Index:>) {
				int length = 0;
				for (int shift = 0; true; shift += 7) {
					byte b = buf[i++];
					length |= (b & 0x7f) << shift;
					if (shift == 28 || b >= 0) break;
				}
				if (length > colferListMax)
					throw new SecurityException(format("colfer: field <:.String:> length %d exceeds %d elements", length, colferListMax));

				<:.TypeNative:>[] a = new <:.TypeNative:>[length];
				for (int ai = 0; ai < length; ai++) {
					<:.TypeNative:> o = new <:.TypeNative:>();
					i = o.unmarshal(buf, i);
					a[ai] = o;
				}
				this.<:.Name:> = a;
				header = buf[i++];
			}
<:else:>
			if (header == (byte) <:.Index:>) {
				this.<:.Name:> = new <:.TypeNative:>();
				i = this.<:.Name:>.unmarshal(buf, i);
				header = buf[i++];
			}
<:end:><:end:>
			if (header != (byte) 0x7f)
				throw new InputMismatchException(format("colfer: unknown header at byte %d", i - 1));
		} catch (IndexOutOfBoundsException e) {
			if (i - offset > colferSizeMax)
				throw new SecurityException(format("colfer: serial exceeds %d bytes", colferSizeMax));
			if (i >= buf.length)
				throw new BufferUnderflowException();
			throw new RuntimeException("colfer: bug", e);
		}

		if (i - offset > colferSizeMax)
			throw new SecurityException(format("colfer: serial exceeds %d bytes", colferSizeMax));
		return i;
	}
<:range .Fields:>
	public <:.TypeNative:><:if .TypeArray:>[]<:end:> get<:.NameTitle:>() {
		return this.<:.Name:>;
	}

	public void set<:.NameTitle:>(<:.TypeNative:><:if .TypeArray:>[]<:end:> value) {
		this.<:.Name:> = value;
	}
<:end:>
	@Override
	public final int hashCode() {
		int h = 1;
<:- range .Fields:>
<:- if eq .Type "bool":>
		h = 31 * h + (this.<:.Name:> ? 1231 : 1237);
<:- else if eq .Type "uint32" "int32":>
		h = 31 * h + this.<:.Name:>;
<:- else if eq .Type "uint64" "int64":>
		h = 31 * h + (int)(this.<:.Name:> ^ this.<:.Name:> >>> 32);
<:- else if eq .Type "float32":>
		h = 31 * h + Float.floatToIntBits(this.<:.Name:>);
<:- else if eq .Type "float64":>
		long _<:.Name:>Bits = Double.doubleToLongBits(this.<:.Name:>);
		h = 31 * h + (int) (_<:.Name:>Bits ^ _<:.Name:>Bits >>> 32);
<:- else if eq .Type "binary":>
		for (byte b : this.<:.Name:>) h = 31 * h + b;
<:- else if .TypeArray:>
		for (<:.TypeNative:> o : this.<:.Name:>) h = 31 * h + (o == null ? 0 : o.hashCode());
<:- else:>
		if (this.<:.Name:> != null) h = 31 * h + this.<:.Name:>.hashCode();
<:- end:><:end:>
		return h;
	}

	@Override
	public final boolean equals(Object o) {
		return o instanceof <:.NameTitle:> && equals((<:.NameTitle:>) o);
	}

	public final boolean equals(<:.NameTitle:> o) {
		return o != null
<:- range .Fields:>
<:- if eq .Type "bool" "uint32" "uint64" "int32" "int64":>
			&& this.<:.Name:> == o.<:.Name:>
<:- else if eq .Type "float32" "float64":>
			&& (this.<:.Name:> == o.<:.Name:> || (this.<:.Name:> != this.<:.Name:> && o.<:.Name:> != o.<:.Name:>))
<:- else if eq .Type "binary":>
			&& java.util.Arrays.equals(this.<:.Name:>, o.<:.Name:>)
<:- else if .TypeArray:>
			&& java.util.Arrays.equals(this.<:.Name:>, o.<:.Name:>)
<:- else:>
			&& java.util.Objects.equals(this.<:.Name:>, o.<:.Name:>)
<:- end:><:end:>;
	}

}
`
