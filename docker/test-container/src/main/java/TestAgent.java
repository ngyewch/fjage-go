import java.util.Map;
import java.util.TreeMap;

import org.arl.fjage.Agent;
import org.arl.fjage.param.ParameterMessageBehavior;

public class TestAgent extends Agent {
  private String[] strings = new String[] {"hello, world"};
  private int[] ints = new int[] {12, -34, 56, -78, 89};
  private float[] floats = new float[] {12.34f, -45.67f, 89f};
  private double[] doubles = new double[] {12.34d, -45.67d, 89d};
  private String string1 = "my string";
  private int int1 = 123;
  private float float1 = 36.47f;
  private double double1 = -25.14d;
  private Map<String, Object> map = new TreeMap<>();

  @Override
  protected void init() {
    super.init();

    map.put("subsystem1", 20.1f);
    map.put("subsystem2", 30.2f);
    map.put("subsystem3", 40.3f);

    add(new ParameterMessageBehavior(TestParam.class));
  }

  public String[] getStrings() {
    return strings;
  }

  public void setStrings(String[] values) {
    strings = values;
  }

  public int[] getInts() {
    return ints;
  }

  public void setInts(int[] values) {
    this.ints = values;
  }

  public float[] getFloats() {
    return floats;
  }

  public void setFloats(float[] values) {
    this.floats = values;
  }

  public double[] getDoubles() {
    return doubles;
  }

  public void setDoubles(double[] values) {
    this.doubles = values;
  }

  public String getString1() {
    return string1;
  }

  public void setString1(String value) {
    this.string1 = value;
  }

  public int getInt1() {
    return int1;
  }

  public void setInt1(int value) {
    this.int1 = value;
  }

  public float getFloat1() {
    return float1;
  }

  public void setFloat1(float value) {
    this.float1 = value;
  }

  public double getDouble1() {
    return double1;
  }

  public void setDouble1(double value) {
    this.double1 = value;
  }

  public Map<String, Object> getMap() {
    return map;
  }

  public void setMap(Map<String, Object> map) {
    this.map = map;
  }
}
