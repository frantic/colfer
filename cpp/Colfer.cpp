#include <cstddef>

#include <chrono>
#include <string>
#include <vector>

// O contains all supported data types.
struct O {
	// b tests booleans.
	bool b;

	// u8 tests unsigned 8-bit integers.
	std::uint8_t u8;

	// u16 tests unsigned 16-bit integers.
	std::uint16_t u16;

	// u32 tests unsigned 32-bit integers.
	std::uint32_t u32;

	// u64 tests unsigned 64-bit integers.
	std::uint64_t u64;

	// u32 tests signed 32-bit integers.
	std::int32_t i32;

	// u64 tests signed 64-bit integers.
	std::int64_t i64;

	// u32 tests 32-bit floating points.
	float f32;

	// u32s tests 32-bit floating point lists.
	std::vector<float> f32s;

	// u64 tests 64-bit floating points.
	double f64;

	// u64s tests 64-bit floating point lists.
	std::vector<double> f64s;

	// t tests timestamps.
	std::chrono::nanoseconds t;

	// a tests binaries.
	std::vector<uint8_t> a;

	// as tests binary lists.
	std::vector<std::vector<std::uint8_t> > as;

	// s tests text.
	std::string s;

	// ss tests text lists.
	std::vector<std::string> ss;

	// o tests nested data structures.
	O* o;

	// os tests data structure lists.
	std::vector<O> os;

	// marshal_len returns the Colfer serial byte size.
	size_t marshal_len();

	// marshal_to encodes O as Colfer into buf and returns the number of bytes written.
	size_t marshal_to(void* buf);

	// unmarshal decodes data as Colfer and returns the number of bytes read.
	size_t umarshal(void* data, size_t len);
};

int main() {
}
